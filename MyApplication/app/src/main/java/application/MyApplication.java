package application;

import android.app.Application;
import android.content.Context;
import android.widget.Toast;

import com.example.wuxiangan.bangbang.R;

import socket.Socket;
import util.ToastUtil;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class MyApplication extends Application {
    private static MyApplication mInstance=null;
    private static Socket socket = new Socket();

    private void initSocket() {
        if (!socket.connect(Config.SERVER_IP, Config.SERVER_PORT)) {
            Toast.makeText(this, R.string.network_error, Toast.LENGTH_SHORT);
        }
    }
    @Override
    public void onCreate() {
        super.onCreate();
        mInstance = this;

        initSocket();
    }

    public static MyApplication getInstance() {
        return mInstance;
    }

    public static Context getContext() {
        return mInstance.getApplicationContext();
    }

    public static Socket getSocket() {return socket;}



}
