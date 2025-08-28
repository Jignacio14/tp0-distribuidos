import logging
import socket


DELIMITER = '\000'


class ServerProtocol:

    def __init__(self, client_ckt):
        self._client_skt = client_ckt

    def __receive_until_delimiter(self):
        bytes = bytearray()
        delimiter = DELIMITER.encode('utf-8')
        found_end = False 
        index = 0

        while not found_end: 
            chunk = self._client_skt.recv(1024)
            if not chunk:
                raise OSError("Connection closed by the client")
            bytes.extend(chunk)

            index = bytes.find(delimiter)
            if index != -1:
                found_end = True
        
        return bytes[:index].decode('utf-8')

            

    def receive_client_info(self):
        try:
            return self.__receive_until_delimiter()
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return None 
        
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