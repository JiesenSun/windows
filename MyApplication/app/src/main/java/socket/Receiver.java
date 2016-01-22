package socket;

import android.os.Handler;

import application.MyApplication;
import go.client.Client;

/**
 * Created by wuxiangan on 2016/1/21.
 */
public abstract class Receiver extends Client.Receiver.Stub{
    //private static Handler handler= null;
    private static Handler handler = new Handler(MyApplication.getContext().getMainLooper());
    // 设置默认处理handler
    public static void setHandler(Handler handler) {
        Receiver.handler = handler;
    }
    // 若不想再默认handler中处理回调，则重写此方法
    public Handler getHandler() { return handler; }
    @Override
    public void Run(final byte[] data) {
        Handler handler = getHandler();
        if (handler != null) {    // 在指定线程中执行回调
            handler.post(new Runnable() {
                @Override
                public void run() {
                    handle(data);
                }
            });
        } else {
            handle(data);         // 没有设置，在当前线程处理回调
        }
    }

    public abstract void handle(byte[] data);
}
