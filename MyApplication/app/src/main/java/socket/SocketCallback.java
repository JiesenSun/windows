package socket;

import android.os.Handler;

import protocol.BaseData;
import socket.BaseSocketCallback;

/**
 * Created by wuxiangan on 2016/1/5.
 */
public class SocketCallback implements BaseSocketCallback {
    @Override
    public Handler getCallbackHandler() {
        return null;
    }

    @Override
    public void onRecvSuccess(BaseData baseData) {
        System.out.println("onRecvSuccess");
    }

    @Override
    public void onRecvFailed(BaseData baseData, String errMsg) {
        System.out.println("onRecvFailed command:" + baseData.command + "  error info:" + errMsg);
    }

    @Override
    public void onSendSuccess(BaseData baseData) {
        System.out.println("onSendSuccess");
    }

    @Override
    public void onSendFailed(BaseData baseData, String errMsg) {
        System.out.println("onSendFailed command:" + baseData.command + "errror info:" + errMsg);
    }
}
