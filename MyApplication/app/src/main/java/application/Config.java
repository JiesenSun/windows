package application;

import util.SpUtil;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class Config {
    public final static String SERVER_IP = "192.168.20.51";
    public final static int SERVER_PORT = 9100;

    private static String username=null;
    private static String password=null;
    private final static String usernameKey="username";
    private final static String passwordKey="password";
    public static String getUsername() {
        if (username == null) {
            username = SpUtil.getSharedPreference().getString(usernameKey,"");
        }
        return username;
    }
    public static String getPassword() {
        if (password == null) {
            password = SpUtil.getSharedPreference().getString(passwordKey,"");
        }
        return password;
    }
}
