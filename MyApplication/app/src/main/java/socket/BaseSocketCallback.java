package socket;

import android.os.Handler;
import protocol.BaseData;

/**
 * Created by wuxiangan on 2015/12/31.
 */
public interface BaseSocketCallback {
    // 获得回调处理handler
    Handler getCallbackHandler();
    // 接收数据成功回调
    void onRecvSuccess(BaseData baseData);
    // 接收数据失败回调  主要超时未接收等
    void onRecvFailed(BaseData baseData, String errMsg);
    // 发送数据成功回调
    void onSendSuccess(BaseData baseData);
    // 发送数据失败回调
    void onSendFailed(BaseData baseData, String errMsg);
}
