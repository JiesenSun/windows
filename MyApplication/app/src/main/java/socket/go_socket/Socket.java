package socket.go_socket;

import java.util.concurrent.ConcurrentHashMap;
import go.client.Client;

/**
 * Created by wuxiangan on 2016/1/21.
 */
public class Socket {
    private Client.Socket socket;
    //private static ConcurrentHashMap<Long, Class<?>> typeMap = new ConcurrentHashMap<Long,Class<?>>();
    public Socket() {
        socket = Client.NewSocket();
    }

    public boolean connect(String ip, long port) {
        return socket.Connect(ip, port);
    }

    public void send(long cmd, byte[] data,Receiver receiver) {
        socket.Send(cmd, data, receiver);
    }

    public void close() {
        socket.Close();
    }
}
