package util;

import android.content.Context;
import android.content.SharedPreferences;
import android.text.TextUtils;

import com.example.wuxiangan.xxx.MyApplication;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by wuxiangan on 2015/12/24.
 */
public class SpUtil {
    public String mSeparator = "$";
    public final String mDefaultSPName = "default_shared_preferences";

    public Context context = MyApplication.getInstance().getApplicationContext();
    public static SharedPreferences mSPInstance = null;
    public static SpUtil mInstance = null;
    public SpUtil() {
        mSPInstance = context.getSharedPreferences(mDefaultSPName, Context.MODE_PRIVATE);
    }

    public static SpUtil getInstance() {
        if (mInstance == null) {
            mInstance = new SpUtil();
        }
        return mInstance;
    }

    public static SharedPreferences getSharedPreference() {
        return mSPInstance;
    }

    public void remove(String key) {
        mSPInstance.edit().remove(key).commit();
    }

    public List<String> getStrings(String key) {
        String value = mSPInstance.getString(key, "");
        List<String> list = new ArrayList<>();
        if (TextUtils.isEmpty(value)) {
            return list;
        }

        String[] ss = value.split("\\"+mSeparator);
        for (String s : ss) {
            list.add(s);
        }
        return list;
    }

    public void putStrings(String key, String value) {
        String oldValue = mSPInstance.getString(key, "");
        if (false == TextUtils.isEmpty(oldValue)) {
            value = oldValue + mSeparator + value;
        }
        SharedPreferences.Editor editor = mSPInstance.edit();
        editor.putString(key, value);
        editor.commit();
    }
}
