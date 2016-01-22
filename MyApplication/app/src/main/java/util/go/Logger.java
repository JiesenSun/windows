package util.go;

import android.util.Log;

import go.client.Client;
import socket.Socket;

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
    private static int level = 0;
    private static Socket socket;
    private static String identify;

    private static void setLevel(int level) {Logger.level = level; }

    public static boolean connect(String ip, long port) {
        if (socket != null) {
            socket.close();
        }

        socket = new Socket();
        return socket.connect(ip, port);
    }
    public static void close() {
        socket.close();
        socket = null;
    }
    public static void send(String tag,int level, String content) {
        if (Logger.level > level) {
            return;
        }

        if (socket != null) {
            //socket.send(0, null, null);
        }
        Log.d(tag, content);
    }
    @Override
    public void Debug(String s) { send(TAG, LOG_LEVEL_DEBUG, s);  }

    @Override
    public void Fatal(String s) {
        send(TAG, LOG_LEVEL_FATAL , s);
    }

    @Override
    public void Info(String s) {
        send(TAG,LOG_LEVEL_INFO , s);
    }

    @Override
    public void Warn(String s) { send(TAG, LOG_LEVEL_WARN , s); }

    public static void d(String tag, String msg) { send(tag, LOG_LEVEL_DEBUG, msg);}
    public static void i(String tag, String msg) { send(tag, LOG_LEVEL_INFO, msg);}
    public static void w(String tag, String msg) { send(tag, LOG_LEVEL_WARN, msg);}
    public static void f(String tag, String msg) { send(tag, LOG_LEVEL_FATAL, msg);}
}
