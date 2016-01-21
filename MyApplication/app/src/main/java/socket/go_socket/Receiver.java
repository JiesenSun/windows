package socket.go_socket;

import android.os.Handler;

import go.client.Client;

/**
 * Created by wuxiangan on 2016/1/21.
 */
public abstract class Receiver extends Client.Receiver.Stub{
    private Handler handler=null;
    public void setHandler(Handler handler) {
        this.handler = handler;
    }
    @Override
    public void Run(final byte[] data) {
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
