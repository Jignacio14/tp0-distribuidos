import logging
import socket


OP_CODE_LEN = 1
BET_LEN = 4


CLIENT_MSG_CODE = 1
OK_CODE = 2
NO_CODE = 3

class ServerProtocol:

    def __init__(self, client_ckt):
        self._client_skt = client_ckt

    
    def receive_bet(self):
        try:
            op_code = self.__receive_op_code()
            if op_code != CLIENT_MSG_CODE:
                logging.error(f"action: receive_message | result: fail | error: invalid op code received: {op_code}")
                return None
            bet_len = self.__receive_bet_lenght()
            return self.__receive_bet(bet_len)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return None 
        
    def __receive_op_code(self) -> int:
        """
        receive the op code (1 byte int) from the socket.
        """
        op_code_byte = self.__receive_all(OP_CODE_LEN)
        return int.from_bytes(op_code_byte, byteorder='big')
    
    def __receive_bet_lenght(self) -> int:
        """
        receive the length of the bet (4 bytes int) from the socket.
        """
        length_bytes = self.__receive_all(BET_LEN)
        return int.from_bytes(length_bytes, byteorder='big')
    
    def __receive_bet(self, length) -> str:
        """"
        Receive a bet of `length` bytes from the socket.
        """
        batch_bytes = self.__receive_all(length)
        return batch_bytes.decode('utf-8')

    def send_confirmation(self, confirmation: bool):
        try:
            op_code = OK_CODE if confirmation else NO_CODE
            self._client_skt.sendall(op_code.to_bytes(1, byteorder='big'))
        except OSError as e:
            logging.error(f"action: send_end_of_batches | result: fail | error: {e}")
            return

    def __receive_all(self, lenght: int) -> bytes:
        """
        Receive exactly `lenght` bytes from the socket.
        """
        data_bytes = bytearray()
        while len(data_bytes) < lenght:
            chunk = self._client_skt.recv(lenght - len(data_bytes))
            if not chunk:
                raise OSError("Connection closed by the client")
            data_bytes.extend(chunk)
        return bytes(data_bytes)
    
    def shutdown(self):
        try:
            if not self._client_skt:
                return
            self._client_skt.shutdown(socket.SHUT_RDWR)
            self._client_skt.close()
        except OSError as e:
            logging.error(f"action: shutdown_connection | result: fail | error: {e}")