package socket.java_socket;

import android.os.Handler;

/**
 * Created by wuxiangan on 2016/1/5.
 */
public class Socket {
    public static BaseSocket mSocket = new BaseSocket();

    public static boolean connect(String ip, int port) {
        return mSocket.connect(ip, port);
    }
    public static void setCallbackHandler(Handler handler) { mSocket.setCallbackHandler(handler);}
    public static void send(int cmd, BaseData baseData, BaseSocketCallback callback) {
        mSocket.send(cmd, baseData, callback);
    }


}
