package application;

import util.SpUtil;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class Config {
    public final static String SERVER_IP = "192.168.20.25";
    public final static int SERVER_PORT = 8888;

    private static String username=null;
    private static String password=null;
    private final static String usernameKey="username";
    private final static String passwordKey="password";
    private static long userID;
    public static long getUserID() { return userID;}
    public static void setUserID(long userID) {  Config.userID = userID; }

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
    public static void setUsername(String username) {
        Config.username = username;
        SpUtil.getSharedPreference().edit().putString(usernameKey, username).commit();
    }
    public static void setPassword(String password) {
        Config.password = password;
        SpUtil.getSharedPreference().edit().putString(passwordKey,password).commit();
    }
}
