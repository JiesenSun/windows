package com.example.wuxiangan.xxx;

import android.animation.Animator;
import android.animation.AnimatorListenerAdapter;
import android.annotation.TargetApi;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.graphics.Color;
import android.support.annotation.NonNull;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AppCompatActivity;
import android.app.LoaderManager.LoaderCallbacks;

import android.content.CursorLoader;
import android.content.Loader;
import android.database.Cursor;
import android.net.Uri;
import android.os.AsyncTask;

import android.os.Build;
import android.os.Bundle;
import android.provider.ContactsContract;
import android.text.TextUtils;
import android.text.method.KeyListener;
import android.util.Log;
import android.view.KeyEvent;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.inputmethod.EditorInfo;
import android.widget.ArrayAdapter;
import android.widget.AutoCompleteTextView;
import android.widget.Button;
import android.widget.CheckBox;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import java.util.ArrayList;
import java.util.List;

import util.SpUtil;

import static android.Manifest.permission.READ_CONTACTS;

/**
 * A login screen that offers login via Phone/password.
 */
public class LoginActivity extends AppCompatActivity {

    /**
     * Id to identity READ_CONTACTS permission request.
     */
    private static final int REQUEST_READ_CONTACTS = 0;

    /**
     * Keep track of the login task to ensure we can cancel it if requested.
     */
    private UserLoginTask mAuthTask = null;

    // UI references.
    private AutoCompleteTextView mPhoneView;
    private EditText mPasswordView;
    private View mProgressView;
    private View mLoginFormView;
    private static final String mSpPhoneKey = "phone_number";
    private static final String mSpLastPhoneKey = "last_phone_number";
    private static final String mSpLastPasswordKey = "last_password";
    private static final String mSpRememberPsdKey = "remember_password";
    private static final String mSpAutoLoginKey = "auto_login";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);

        initView();
    }
    private void initView() {
        // 加载资源
        final CheckBox rememberPsdView = (CheckBox)findViewById(R.id.remember_psd);
        final CheckBox autoLoginView = (CheckBox)findViewById(R.id.auto_login);
        mPhoneView = (AutoCompleteTextView) findViewById(R.id.phone);
        mPasswordView = (EditText) findViewById(R.id.password);
        mLoginFormView = findViewById(R.id.login_form);
        mProgressView = findViewById(R.id.login_progress);
        final Button mPhoneSignInButton = (Button) findViewById(R.id.phone_sign_in_button);
        // 设置监听
        rememberPsdView.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                SpUtil.getSharedPreference().edit().putBoolean(mSpRememberPsdKey, rememberPsdView.isChecked()).commit();
            }
        });
        autoLoginView.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                SpUtil.getSharedPreference().edit().putBoolean(mSpAutoLoginKey, autoLoginView.isChecked()).commit();
            }
        });
        mPhoneSignInButton.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View view) {
                mPhoneSignInButton.setBackgroundColor(Color.DKGRAY);
                attemptLogin();
            }
        });

        mPasswordView.setOnEditorActionListener(new TextView.OnEditorActionListener() {
            @Override
            public boolean onEditorAction(TextView textView, int id, KeyEvent keyEvent) {
                if (id == R.id.login || id == EditorInfo.IME_NULL) {
                    attemptLogin();
                    return true;
                }
                return false;
            }
        });

        // 初始化界面数据
        // Set up the login form
        SharedPreferences sp = SpUtil.getSharedPreference();
        boolean rememberPsd = sp.getBoolean(mSpRememberPsdKey, false);
        boolean autoLogin = sp.getBoolean(mSpAutoLoginKey, false);
        String lastPhoneNumber = sp.getString(mSpLastPhoneKey, "");
        String lastPassword = sp.getString(mSpLastPasswordKey, "");
        List<String> phonelist = SpUtil.getInstance().getStrings(mSpPhoneKey);
        // 设置checkbox
        rememberPsdView.setChecked(rememberPsd);
        autoLoginView.setChecked(autoLogin);
        // 自动登录
        if (rememberPsd && autoLogin && !TextUtils.isEmpty(lastPhoneNumber) && !TextUtils.isEmpty(lastPassword)) {
            attemptLogin();
            return;
        }
        // 设置账号匹配
        ArrayAdapter<String> adapter = new ArrayAdapter<String>(LoginActivity.this, R.layout.support_simple_spinner_dropdown_item, phonelist);
        mPhoneView.setAdapter(adapter);
        mPhoneView.setText(lastPhoneNumber);
        // 若记住密码，则登录设置密码
        if (rememberPsd) {
            mPasswordView.setText(lastPassword);
        }
        // 已存在账户，聚焦密码
        if (false == TextUtils.isEmpty(lastPhoneNumber)) {
            if (!rememberPsd || TextUtils.isEmpty(lastPassword)) {
                mPasswordView.requestFocus();
            }
        }
    }
    /**
     * Attempts to sign in or register the account specified by the login form.
     * If there are form errors (invalid Phone, missing fields, etc.), the
     * errors are presented and no actual login attempt is made.
     */
    private void attemptLogin() {
        if (mAuthTask != null) {
            return;
        }

        // Reset errors.
        mPhoneView.setError(null);
        mPasswordView.setError(null);

        // Store values at the time of the login attempt.
        String Phone = mPhoneView.getText().toString();
        String password = mPasswordView.getText().toString();

        boolean cancel = false;
        View focusView = null;

        // Check for a valid password, if the user entered one.
        if (!TextUtils.isEmpty(password) && !isPasswordValid(password)) {
            mPasswordView.setError(getString(R.string.error_invalid_password));
            focusView = mPasswordView;
            cancel = true;
        }

        // Check for a valid Phone address.
        if (TextUtils.isEmpty(Phone)) {
            mPhoneView.setError(getString(R.string.error_field_required));
            focusView = mPhoneView;
            cancel = true;
        } else if (!isPhoneValid(Phone)) {
            mPhoneView.setError(getString(R.string.error_invalid_phone));
            focusView = mPhoneView;
            cancel = true;
        }

        if (cancel) {
            // There was an error; don't attempt login and focus the first
            // form field with an error.
            focusView.requestFocus();
        } else {
            // Show a progress spinner, and kick off a background task to
            // perform the user login attempt.
            showProgress(true);
            mAuthTask = new UserLoginTask(Phone, password);
            mAuthTask.execute((Void) null);
        }
    }

    private boolean isPhoneValid(String Phone) {
        //TODO: Replace this with your own logic
        return Phone.length() == 11;
    }

    private boolean isPasswordValid(String password) {
        //TODO: Replace this with your own logic
        return password.length() > 4;
    }

    /**
     * Shows the progress UI and hides the login form.
     */
    @TargetApi(Build.VERSION_CODES.HONEYCOMB_MR2)
    private void showProgress(final boolean show) {
        // On Honeycomb MR2 we have the ViewPropertyAnimator APIs, which allow
        // for very easy animations. If available, use these APIs to fade-in
        // the progress spinner.
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.HONEYCOMB_MR2) {
            int shortAnimTime = getResources().getInteger(android.R.integer.config_shortAnimTime);

            mLoginFormView.setVisibility(show ? View.GONE : View.VISIBLE);
            mLoginFormView.animate().setDuration(shortAnimTime).alpha(
                    show ? 0 : 1).setListener(new AnimatorListenerAdapter() {
                @Override
                public void onAnimationEnd(Animator animation) {
                    mLoginFormView.setVisibility(show ? View.GONE : View.VISIBLE);
                }
            });

            mProgressView.setVisibility(show ? View.VISIBLE : View.GONE);
            mProgressView.animate().setDuration(shortAnimTime).alpha(
                    show ? 1 : 0).setListener(new AnimatorListenerAdapter() {
                @Override
                public void onAnimationEnd(Animator animation) {
                    mProgressView.setVisibility(show ? View.VISIBLE : View.GONE);
                }
            });
        } else {
            // The ViewPropertyAnimator APIs are not available, so simply show
            // and hide the relevant UI components.
            mProgressView.setVisibility(show ? View.VISIBLE : View.GONE);
            mLoginFormView.setVisibility(show ? View.GONE : View.VISIBLE);
        }
    }

    /**
     * Represents an asynchronous login/registration task used to authenticate
     * the user.
     */
    public class UserLoginTask extends AsyncTask<Void, Void, Boolean> {

        private final String mPhone;
        private final String mPassword;

        UserLoginTask(String Phone, String password) {
            mPhone = Phone;
            mPassword = password;
        }

        @Override
        protected Boolean doInBackground(Void... params) {
            // TODO: attempt authentication against a network service.

            try {
                // Simulate network access.
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                return false;
            }

            // TODO: register the new account here.
            return true;
        }

        @Override
        protected void onPostExecute(final Boolean success) {
            mAuthTask = null;
            showProgress(false);

            if (success) {
                SpUtil.getInstance().putStrings(mSpPhoneKey, mPhone);
                SpUtil.getSharedPreference().edit().putString(mSpLastPhoneKey, mPhone);
                SpUtil.getSharedPreference().edit().putString(mSpLastPasswordKey, mPassword).commit();
                // 执行页面跳转
                //finish();
                Toast.makeText(MyApplication.getInstance().getApplicationContext(), "等陆成功", Toast.LENGTH_SHORT).show();
            } else {
                mPasswordView.setError(getString(R.string.error_incorrect_password));
                mPasswordView.requestFocus();
                Toast.makeText(MyApplication.getInstance().getApplicationContext(), "等陆失败", Toast.LENGTH_SHORT).show();

            }
        }

        @Override
        protected void onCancelled() {
            mAuthTask = null;
            showProgress(false);
        }
    }
}

