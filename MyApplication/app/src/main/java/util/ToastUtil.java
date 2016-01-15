package util;

import android.content.Context;
import android.widget.Toast;

import application.MyApplication;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class ToastUtil {
    private static Context context = MyApplication.getInstance().getApplicationContext();

    public static void makeText(String text) {
        Toast.makeText(context, text, Toast.LENGTH_SHORT).show();
    }

    public static void makeText(int strID) {
        Toast.makeText(context, strID, Toast.LENGTH_SHORT).show();
    }

    public static void makeText(String text, int duration) {
        Toast.makeText(context, text, duration).show();
    }

    public static void makeText(int strID, int duration) {
        Toast.makeText(context, strID, duration).show();
    }
}
