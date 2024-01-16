import socket
from concurrent.futures import ThreadPoolExecutor

def send_tcp_packet(target_host, target_port, data):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((target_host, target_port))
        s.sendall(data.encode())
        response = s.recv(1024)
        print(f"Received from {target_host}:{target_port}: {response.decode()}")

def main():
    target_ip = '127.0.0.1' 
    target_port = 3090  

    data_to_send = "SUM 1 2 33 10 END"

    with ThreadPoolExecutor(max_workers=1000) as executor:
        futures = [
            executor.submit(send_tcp_packet, target_ip, target_port, data_to_send)
            for _ in range(1000)
        ]

        for future in futures:
            future.result()

if __name__ == "__main__":
    main()