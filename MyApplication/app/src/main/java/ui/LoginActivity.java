package ui;

import android.app.Activity;
import android.os.Bundle;
import android.view.Window;

import com.example.wuxiangan.bangbang.R;

import ui.base.BaseActivity;
import util.ActivityCollector;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class LoginActivity extends BaseActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.login_activity);

        setTitleLeft(R.string.back);
        setTitleRight(R.string.register);
        setTitle(R.string.login);
    }
}
