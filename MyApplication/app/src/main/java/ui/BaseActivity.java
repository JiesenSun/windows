package ui;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;

import util.ActivityCollector;

/**
 * Created by wuxiangan on 2016/1/14.
 */
public class BaseActivity extends Activity {
    private final static String TAG="BaseActivity";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        Log.d(TAG, getClass().getSimpleName());
        ActivityCollector.addActivity(this);
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        ActivityCollector.removeActivity(this);
    }
}
