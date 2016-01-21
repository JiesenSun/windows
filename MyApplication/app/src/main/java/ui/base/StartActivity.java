package ui.base;

import android.content.Intent;
import android.os.Bundle;
import android.os.Message;
import android.text.TextUtils;
import android.os.Handler;

import application.Config;
import socket.java_socket.Socket;
import ui.LoginActivity;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class StartActivity extends BaseActivity {
    private Handler handler = new Handler() {
        @Override
        public void handleMessage(Message msg) {
            super.handleMessage(msg);
        }
    };
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        Socket.setCallbackHandler(handler);

        String username = Config.getUsername();
        String password = Config.getPassword();
        if (!TextUtils.isEmpty(username) && !TextUtils.isEmpty(password)) {
            // LoginReq
        } else {
            Intent intent = new Intent(this, LoginActivity.class);
            startActivity(intent);
        }
    }
}
