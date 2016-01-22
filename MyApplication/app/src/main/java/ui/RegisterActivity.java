package ui;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.example.wuxiangan.bangbang.R;
import com.google.protobuf.InvalidProtocolBufferException;

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
public class RegisterActivity extends BaseActivity {
    private final static String TAG = "registerActivity";
    private final static int REGISTER_STEP_USERNAME = 0;
    private final static int REGISTER_STEP_PASSWORD = 1;
    private static String username;
    private static String password;
    private static int step = 0;
    private EditText registerET;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.register);

        registerET = (EditText)findViewById(R.id.register_et);
        if (step == REGISTER_STEP_USERNAME)
            registerET.setHint(R.string.password_input_tip);
        else if (step == REGISTER_STEP_PASSWORD)
            registerET.setHint(R.string.password_input_tip);

        Button nextBtn = (Button)findViewById(R.id.next_btn);
        nextBtn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (step == REGISTER_STEP_USERNAME) {
                    username = registerET.getText().toString();
                    step = REGISTER_STEP_PASSWORD;
                    Intent intent = new Intent(RegisterActivity.this, RegisterActivity.class);
                    startActivity(intent);
                } else if (step == REGISTER_STEP_PASSWORD) {
                    password = registerET.getText().toString();
                    step = REGISTER_STEP_USERNAME;

                    // 正式注册
                    Protocol.register_request.Builder builder = Protocol.register_request.newBuilder();
                    builder.setUsername(username);
                    builder.setPassword(password);
                    byte[] data = builder.build().toByteArray();
                    MyApplication.getSocket().send(Command.CLINET_CMD_REGISTER,data, registerReceiver);
                }
            }
        });
    }

    private Receiver registerReceiver = new Receiver() {

        @Override
        public void handle(byte[] data) {
            Protocol.register_response register_response = null;
            try {
                register_response = Protocol.register_response.parseFrom(data);
            } catch (InvalidProtocolBufferException e) {
                e.printStackTrace();
                Logger.d(TAG, "server_exception");
                ToastUtil.makeText(R.string.server_exception);
                return;
            }

            int errCode = register_response.getErrorCode();
            if (errCode == 0) {
                //LoginActivity.actionStart(RegisterActivity.this, username, password);
                Config.setUserID(register_response.getUserId());
                Intent intent = new Intent(RegisterActivity.this, MainActivity.class);
                startActivity(intent);
                ActivityCollector.finish(RegisterActivity.class);
            } else if (errCode == Command.CLIENT_ERR_USER_EXIST) {
                ToastUtil.makeText(R.string.user_registered);
            }
        }
    };
}
