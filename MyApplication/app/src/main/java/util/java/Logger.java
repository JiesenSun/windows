package util.java;

/**
 * Created by wuxiangan on 2015/12/30.
 */

import android.util.Log;

import org.json.JSONObject;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.nio.ByteBuffer;

public class Logger {
    private final  static String TAG="Logger";
    public static boolean isDebug = false;
    public static boolean isEnable = false;
    public static String id;
    public static String  ip;
    public static int port;
    private static Socket socket;
    private static DataOutputStream out;

    public static void idendify(String id) {
        Logger.id = id;
    }
    public static void debug(boolean debug) {
        Logger.isDebug = debug;
    }

    public static void enable(boolean enable) {
        Logger.isEnable = enable;
    }

    public static void logAddree(String ip, int port) {
        Logger.ip = ip;
        Logger.port = port;
    }
    public static boolean connect(String ip, int port) {
        if (Logger.isConnected()) {
            return true;
        }
        Logger.ip = ip;
        Logger.port = port;

        try {
            socket = new Socket(ip, port);
            out = new DataOutputStream(socket.getOutputStream());
        } catch (IOException e) {
            Log.i(TAG, "connect server failed");
            socket = null;
            out = null;
            return false;
        }
        Log.d(TAG, "connect server success!!!!");
        return true;
    }

    public static boolean isConnected() {
        if (socket == null)
            return false;

        try {
             socket.sendUrgentData(0);
             socket.sendUrgentData(0);
        } catch (Exception e) {
            Log.i(TAG,"disconnect the server!!!");
            Log.i(TAG, e.toString());
            Logger.close();
            return false;
        }
        return  true;
    }

    public static void close() {
        try {
            out.close();
            socket.close();
        } catch (Exception e) {
            Log.i(TAG, e.toString());
        }
        socket = null;
        out = null;
    }
    private static void write(String tag, String level, String msg) {
        JSONObject jsonObject = new JSONObject();
        String logBody;
        try {
            jsonObject.put("identify", Logger.id);
            jsonObject.put("tag", tag);
            jsonObject.put("level", level);
            jsonObject.put("msg", msg);
            logBody = jsonObject.toString();
        } catch (Exception e) {
            Log.i(TAG, e.toString());
            return ;
        }
        ByteBuffer byteBuffer = ByteBuffer.allocate(2048);
        byteBuffer.putShort((short)logBody.length());
        byteBuffer.put(logBody.getBytes());
        byte[] bytes = new byte[byteBuffer.position()];
        byteBuffer.position(0);
        byteBuffer.get(bytes);
        try {
            out.write(bytes);
        } catch (Exception e) {
            Log.i(TAG, e.toString());
            Logger.close();
        }
    }
    public static void d(String tag, String msg) {
        if (isEnable && tag != null && msg != null){
            if (isDebug && connect(Logger.ip, Logger.port)) {
                write("debug", tag, msg);
            }
            Log.i(tag, msg);
        }
    }
    public static void i(String tag, String msg) {
        if (isEnable && tag != null && msg != null){
            if (isDebug && connect(Logger.ip, Logger.port)) {
                write("info", tag, msg);
            }
            Log.i(tag, msg);
        }
    }

    public static void w(String tag, String msg) {
        if (isEnable && tag != null && msg != null){
            if (isDebug && connect(Logger.ip, Logger.port)) {
                write("warn", tag, msg);
            }
            Log.i(tag, msg);
        }
    }
    public static void e(String tag, String msg) {
        if (isEnable && tag != null && msg != null){
            if (isDebug && connect(Logger.ip, Logger.port)) {
                write("error", tag, msg);
            }
            Log.i(tag, msg);
        }
    }
    public static void f(String tag, String msg) {
        if (isEnable && tag != null && msg != null){
            if (isDebug && connect(Logger.ip, Logger.port)) {
                write("fatal", tag, msg);
            }
            Log.i(tag, msg);
        }
    }
}
