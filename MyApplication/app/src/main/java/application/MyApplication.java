package application;

import android.app.Application;
import android.content.Context;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class MyApplication extends Application {
    private static MyApplication mInstance=null;

    @Override
    public void onCreate() {
        super.onCreate();
        mInstance = this;
    }

    public static MyApplication getInstance() {
        return mInstance;
    }

    public static Context getContext() {
        return mInstance.getApplicationContext();
    }
}
