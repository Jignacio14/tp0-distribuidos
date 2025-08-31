import logging
import socket


END_DELIMITER = '\000'
BATCH_DELIMITER = '\t'


class ServerProtocol:

    def __init__(self, client_ckt):
        self._client_skt = client_ckt
        self._conection_to_end = False

    def __receive_until_delimiter(self):
        bytes = bytearray()
        end_delimiter = END_DELIMITER.encode('utf-8')
        batch_delimiter = BATCH_DELIMITER.encode('utf-8')
        end_delimiter_index = -1
        batch_delimiter_index = -1

        while True: 
            chunk = self._client_skt.recv(1024)
            if not chunk:
                raise OSError("Connection closed by the client")
            bytes.extend(chunk)

            end_delimiter_index = bytes.find(end_delimiter)
            batch_delimiter_index = bytes.find(batch_delimiter)

            if end_delimiter_index != -1:
                self._conection_to_end = True
                return bytes[:end_delimiter_index].decode('utf-8').replace(BATCH_DELIMITER, ''), False
            elif batch_delimiter_index != -1:
                return bytes[:batch_delimiter_index].decode('utf-8').replace(BATCH_DELIMITER, ''), True

    def receive_client_info(self):
        try:
            return self.__receive_until_delimiter()
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return '', False

    def send_confirmation(self, confirmation: bool):
        try:
            msg = 'OK' if confirmation else 'NO'
            self._client_skt.sendall(msg.encode('utf-8'))
        except OSError as e:
            logging.error(f"action: send_message | result: fail | error: {e}")
            return False
        return True
    
    def shutdown(self):
        try:
            self._client_skt.shutdown(socket.SHUT_RDWR)
            self._client_skt.close()
        except OSError as e:
            logging.error(f"action: shutdown_connection | result: fail | error: {e}")
            return False
        return True