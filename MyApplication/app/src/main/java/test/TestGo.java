package test;

import android.test.InstrumentationTestCase;

import java.util.regex.Pattern;


/**
 * Created by wuxiangan on 2016/1/18.
 */
public class TestGo extends InstrumentationTestCase {
    public void testHello() {
        String passwordReg="[\\w]{6,20}";
        Pattern pattern = Pattern.compile(passwordReg);
        boolean matcher = pattern.matcher("wuxiangan").matches();
        System.out.println(matcher);
    }
}
