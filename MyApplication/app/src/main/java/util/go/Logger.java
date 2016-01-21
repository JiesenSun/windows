package util.go;

import android.util.Log;

import go.client.Client;
import socket.go_socket.Socket;

/**
 * Created by wuxiangan on 2016/1/21.
 */
public class Logger extends Client.Logger.Stub {
    private static final String TAG="logger";
    private static final int LOG_LEVEL_DEBUG = 0;
    private static final int LOG_LEVEL_INFO = 1;
    private static final int LOG_LEVEL_WARN = 2;
    private static final int LOG_LEVEL_FATAL = 3;
    private static final int LOG_LEVEL_MAX = 4;
    private int level = 0;
    private Socket socket;
    private String identify;

    private void setLevel(int level) {this.level = level; }

    public boolean connect(String ip, long port) {
        if (socket != null) {
            socket.close();
        }

        socket = new Socket();
        return socket.connect(ip, port);
    }
    public void close() {
        socket.close();
        socket = null;
    }
    public void send(int level, String content) {
        if (this.level > level) {
            return;
        }

        if (socket != null) {
            //socket.send(0, null, null);
        }
        Log.d(TAG, content);
    }
    @Override
    public void Debug(String s) {
        send(LOG_LEVEL_DEBUG, s);
    }

    @Override
    public void Fatal(String s) {
        send(LOG_LEVEL_INFO, s);
    }

    @Override
    public void Info(String s) {
        send(LOG_LEVEL_WARN, s);
    }

    @Override
    public void Warn(String s) {
        send(LOG_LEVEL_FATAL, s);
    }
}
