import logging
import socket

from common.utils import Bet

DELIMITER = '\n'


class ServerProtocol:

    def __init__(self, client_ckt):
        self._client_skt = client_ckt
        self

    def __receive_until_delimiter(self):
        data = b''
        while True:
            chunk = self._client_skt.recv(1024)
            if not chunk:
                break
            data += chunk
            if data.endswith(DELIMITER.encode()):
                break
        return data.decode('utf-8').rstrip(DELIMITER)
    
    def receive_client_info(self):
        try:
            msg = self.__receive_until_delimiter()
            msg_parts = msg.split(',')
            if len(msg_parts) != 6:
                return None
            
            return Bet(msg_parts[0], msg_parts[1], msg_parts[2], msg_parts[3], msg_parts[4], msg_parts[5])
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