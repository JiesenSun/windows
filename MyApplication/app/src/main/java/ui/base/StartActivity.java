package ui.base;

import android.content.Intent;
import android.os.Bundle;
import android.os.Message;
import android.text.TextUtils;
import android.os.Handler;

import application.Config;
import socket.Receiver;
import ui.LoginActivity;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class StartActivity extends BaseActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        //Receiver.setHandler(new Handler(this.getMainLooper()));
        String username = Config.getUsername();
        String password = Config.getPassword();
        if (!TextUtils.isEmpty(username) && !TextUtils.isEmpty(password)) {
            LoginActivity.actionStart(this, username, password);
        } else {
            Intent intent = new Intent(this, LoginActivity.class);
            startActivity(intent);
        }
        finish();
    }
}
