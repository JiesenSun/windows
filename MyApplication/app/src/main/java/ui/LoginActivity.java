package ui;

import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.example.wuxiangan.bangbang.R;
import com.google.protobuf.InvalidProtocolBufferException;

import java.util.regex.Pattern;

import application.Config;
import application.MyApplication;
import protocol.Protocol;
import socket.Command;
import socket.Receiver;
import ui.base.BaseActivity;
import util.ActivityCollector;
import util.ToastUtil;
import util.go.Logger;


/**
 * Created by wuxiangan on 2016/1/15.
 */
public class LoginActivity extends BaseActivity {
    private final static String TAG="loginActivity";
    private final static String passwordReg="^[\\w]{6,20}$";
    private static Pattern pattern = Pattern.compile(passwordReg);
    private static String username;
    private static String password;
    private Button loginBtn;
    // 启动该活动方法
    public static void actionStart(Context context, String username, String password) {
        Intent intent = new Intent(context, LoginActivity.class);
        intent.putExtra("username", username);
        intent.putExtra("password", password);
        context.startActivity(intent);
    }
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.login_activity);

        setTitleLeft(R.string.back);
        setTitleRight(R.string.register);
        setTitle(R.string.login);

        final EditText usernameET = (EditText)findViewById(R.id.username_et);
        final EditText passwordET = (EditText)findViewById(R.id.password_et);
        loginBtn = (Button)findViewById(R.id.login_btn);
        loginBtn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                username = usernameET.getText().toString();
                password = passwordET.getText().toString();
                LoginActivity.this.login(username, password);
            }
        });

        Intent intent = getIntent();
        username = intent.getStringExtra("username");
        password = intent.getStringExtra("password");
        if (!TextUtils.isEmpty(username) && !TextUtils.isEmpty(password)) {
            login(username, password);
        }
    }

    private void login(String username, String password) {
        if (username.length() != 11) {
            ToastUtil.makeText(R.string.phonenumber_fmt_error);
            return;
        }
        if (!pattern.matcher(password).matches()) {
            ToastUtil.makeText(R.string.password_fmt_error);
            return;
        }
        // send login request
        Protocol.login_request.Builder builder = Protocol.login_request.newBuilder();
        builder.setUsername(username);
        builder.setPassword(password);
        MyApplication.getSocket().send(Command.CLINET_CMD_LOGIN, builder.build().toByteArray(),loginReceiver);

    }

    private Receiver loginReceiver = new Receiver() {
        @Override
        public void handle(byte[] data) {
            Protocol.login_response login_response;
            try {
                login_response = Protocol.login_response.parseFrom(data);
            } catch (InvalidProtocolBufferException e) {
                e.printStackTrace();
                return;
            }
            if (Command.CLIENT_ERR_NO_ERROR == login_response.getErrorCode()) {
                Config.setUserID(login_response.getUserId());
                Config.setUsername(username);
                Config.setPassword(password);
                // startActivity
                Intent intent = new Intent(LoginActivity.this, MainActivity.class);
                startActivity(intent);
                // finish activity
                ActivityCollector.finish(LoginActivity.class);
            } else {
                Logger.d(TAG,"login error code:"+login_response.getErrorCode());
                ToastUtil.makeText(R.string.username_password_error);
            }
        }
    };
}
