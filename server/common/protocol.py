import logging
import socket

OP_CODE_LEN = 1
BATCH_LEN = 4

BATCH_SEND_OP_CODE = 1
BATCH_RECEIVED_OK_CODE = 2
BATCH_RECEIVED_FAIL_CODE = 3
END_OF_COMMUNICATION_CODE = 4

class ServerProtocol:

    def __init__(self, client_ckt):
        self._client_skt = client_ckt

    def receive_batch(self) -> str:
        try:
            op_code = self.__receive_op_code()
            if op_code == END_OF_COMMUNICATION_CODE:
                return False, ''
            length = self.__receive_batch_lenght()
            client_info = self.__receive_batch(length)
            return True, client_info
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return ''

    def send_confirmation(self, confirmation: bool):
        try:
            msg = 'OK' if confirmation else 'NO'
            self._client_skt.sendall(msg.encode('utf-8'))
        except OSError as e:
            logging.error(f"action: send_message | result: fail | error: {e}")
            return False
        return True
    
    def __receive_op_code(self) -> int:
        op_code_byte = self.__receive_all(OP_CODE_LEN)
        op_code = int.from_bytes(op_code_byte, byteorder='big')
        return op_code
        
    def __receive_batch_lenght(self) -> int:
        length_bytes = self.__receive_all(BATCH_LEN)
        length = int.from_bytes(length_bytes, byteorder='big')
        return length

    def __receive_batch(self, length) -> str:
        batch_bytes = self.__receive_all(length)
        return batch_bytes.decode('utf-8')

    def send_bad_bets(self, count: int):
        self.__send_template(BATCH_RECEIVED_FAIL_CODE, count)

    def send_batches_received_successfully(self, bets_count: int):
        self.__send_template(BATCH_RECEIVED_OK_CODE, bets_count)

    def __send_template(self, code: int, number: int):
        try:
            self._client_skt.sendall(code.to_bytes(1, byteorder='big'))
            self._client_skt.sendall(number.to_bytes(4, byteorder='big'))
        except OSError as e:
            logging.error(f"action: send_template | result: fail | error: {e}")
            return

    def send_end_of_batches(self):
        try:
            self._client_skt.sendall(END_OF_COMMUNICATION_CODE.to_bytes(1, byteorder='big'))
        except OSError as e:
            logging.error(f"action: send_end_of_batches | result: fail | error: {e}")
            return

    def __receive_all(self, lenght: int) -> bytes:
        data_bytes = bytearray()
        while len(data_bytes) < lenght:
            chunk = self._client_skt.recv(lenght - len(data_bytes))
            if not chunk:
                raise OSError("Connection closed by the client")
            data_bytes.extend(chunk)
        return bytes(data_bytes)

    def shutdown(self):
        try:
            self._client_skt.shutdown(socket.SHUT_RDWR)
            self._client_skt.close()
        except OSError as e:
            logging.error(f"action: shutdown_connection | result: fail | error: {e}")
            return False
        return True