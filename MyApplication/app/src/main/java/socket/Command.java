package socket;

/**
 * Created by wuxiangan on 2016/1/22.
 */
public class Command {
    // command
    public final static long CLIENT_CMD_HEARTBEAT = 10000;
    public final static long CLINET_CMD_REGISTER = 10001;
    public final static long CLINET_CMD_LOGIN = 10002;

    // error code
    public final static long CLIENT_ERR_NO_ERROR = 0;
    public final static long CLIENT_ERR_UNKNOW_ERROR = 1;
    public final static long CLIENT_ERR_USER_NOT_EXIST = 2;
    public final static long CLIENT_ERR_PASSWORD_ERROR = 3;
    public final static long CLIENT_ERR_LOGIN_CONFLICT = 4;
    public final static long CLIENT_ERR_USER_EXIST = 5;
}
