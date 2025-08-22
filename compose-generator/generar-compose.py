import sys

YAML_BODY = """"
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net

  client1:
    container_name: client1
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=1
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

def generate_script(file_destination, num_clients):
    print("Client number:", num_clients)
    print("File destination:", file_destination)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("You should call this script as: ./generar-compose.py <output_file> <num_clients>")
        sys.exit(1)
    
    file_destination = sys.argv[1]
    client_nums = int(sys.argv[2])
    generate_script(file_destination, client_nums)

