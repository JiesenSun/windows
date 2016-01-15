package ui;

import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

import com.example.wuxiangan.bangbang.R;

import java.util.regex.Pattern;

import ui.base.BaseActivity;
import util.ToastUtil;


/**
 * Created by wuxiangan on 2016/1/15.
 */
public class LoginActivity extends BaseActivity {
    private final static String passwordReg="^[\\w]{6,20}$'";
    private static Pattern pattern = Pattern.compile(passwordReg);
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.login_activity);

        setTitleLeft(R.string.back);
        setTitleRight(R.string.register);
        setTitle(R.string.login);

        final EditText usernameET = (EditText)findViewById(R.id.username_et);
        final EditText passwordET = (EditText)findViewById(R.id.password_et);
        Button loginBtn = (Button)findViewById(R.id.login_btn);
        loginBtn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String username = usernameET.getText().toString();
                String password = passwordET.getText().toString();

                if (username.length() != 11) {
                    ToastUtil.makeText(R.string.phonenumber_fmt_error);
                    return;
                }
                if (!pattern.matcher(password).matches()) {
                    ToastUtil.makeText(R.string.password_fmt_error);
                    return ;
                }
            }
        });
    }
}
