package test;

import android.content.SharedPreferences;
import android.test.InstrumentationTestCase;
import android.util.Log;

import java.util.List;

import util.SpUtil;

/**
 * Created by wuxiangan on 2015/12/24.
 */
public class SpUtilTest extends InstrumentationTestCase {
    private String mTestKey = "test_key";
    public void testRemove() {

        Log.d(this.getName(), "remove");
    }

    public void testStrings() {
        SpUtil.getInstance().remove(mTestKey);
        List<String> phonelist = SpUtil.getInstance().getStrings(mTestKey);
        Log.d(this.getName(), phonelist.toString());
        assertEquals(phonelist.size(), 0);

        SpUtil.getInstance().putStrings(mTestKey, "hello world");
        phonelist = SpUtil.getInstance().getStrings(mTestKey);
        Log.d(this.getName(), phonelist.toString());
        assertEquals(phonelist.size(),1);

        SpUtil.getInstance().putStrings(mTestKey, "hello world");
        phonelist = SpUtil.getInstance().getStrings(mTestKey);
        Log.d(this.getName(), phonelist.toString());
        assertEquals(phonelist.size(),2);
    }
}
