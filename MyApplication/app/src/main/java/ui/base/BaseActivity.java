package ui.base;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.DialogInterface;
import android.os.Bundle;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.Window;
import android.widget.FrameLayout;

import com.example.wuxiangan.bangbang.R;

import util.ActivityCollector;

/**
 * Created by wuxiangan on 2016/1/14.
 */
public class BaseActivity extends Activity {
    private final static String TAG="BaseActivity";
    private TitleLayout titleLayout = null;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        Log.d(TAG, getClass().getSimpleName());
        requestWindowFeature(Window.FEATURE_NO_TITLE);

        ActivityCollector.addActivity(this);
    }
    @Override
    protected void onDestroy() {
        super.onDestroy();
        ActivityCollector.removeActivity(this);
    }
    @Override
    public void setContentView(int layoutResID) {
        super.setContentView(R.layout.base_activity);
        FrameLayout mainlayout = (FrameLayout)findViewById(R.id.main_layout);
        LayoutInflater.from(this).inflate(layoutResID, mainlayout);

        titleLayout = (TitleLayout)findViewById(R.id.title_layout);
    }

    public void setTitleLeft(CharSequence text) {
        titleLayout.setTitleLeft(text);
    }
    public void setTitleLeft(int resID) {
        titleLayout.setTitleLeft(getText(resID));
    }

    public void setTitleRight(CharSequence text) {
        titleLayout.setTitleRight(text);
    }
    public void setTitleRight(int resID) {
        titleLayout.setTitleRight(getText(resID));
    }

    @Override
    public void setTitle(CharSequence title) {
        titleLayout.setTitle(title);
    }
    @Override
    public void setTitle(int resID) {
        titleLayout.setTitle(getText(resID));
    }


}
