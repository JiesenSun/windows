package socket.java_socket;

import android.os.Handler;
import android.os.HandlerThread;

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.util.LinkedList;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantLock;

import client.BeatHeart;
import util.java.Logger;

/**
 * Created by wuxiangan on 2015/12/30.
 */
public class BaseSocket {
    private static final String TAG = "BaseSocket";
    private static final int mSecond = 1000;
    private String mIP;
    private int mPort;
    private DataInputStream in;
    private DataOutputStream out;
    private Socket socket;
    private boolean mConnected = false;
    private static final Object mConnectedLock = new Object();
    private int mSendTimeout = 3000;
    private int mRecvTimeout = 10000;
    private boolean mExit = true;

    private ReentrantLock sendTaskLock = new ReentrantLock();
    private Object sendTaskLockObject = new Object();
    private LinkedList<BaseData> sendTaskList = new LinkedList<BaseData>();
    private ConcurrentHashMap<Integer, BaseSocketCallback> mCallbackMap = null;
    private ConcurrentHashMap<Integer, BaseSocketCallback> mPushCallbackMap = null;

    private Handler mTimeoutHandler = null;
    private Handler mCallbackHandler = null;
    private HandlerThread mTimeoutThread = null;

    // 设置回调处理Handler
    public void setCallbackHandler(Handler handler) {mCallbackHandler = handler;}

    // 接收数据成功回调
    public void onRecvSuccess(final BaseSocketCallback callback, final BaseData baseData) {
        Handler handler = callback.getCallbackHandler();
        if (handler == null) {
            handler = mCallbackHandler;
        }
        if (handler != null) {
            handler.post(new Runnable() {
                @Override
                public void run() {
                    callback.onRecvSuccess(baseData);
                }
            });
        } else {
            callback.onRecvSuccess(baseData);
        }
    }

    // 接收数据失败回调  主要超时未接收等
    public void onRecvFailed(final BaseSocketCallback callback, final BaseData baseData, final String errMsg) {
        Handler handler = callback.getCallbackHandler();
        if (handler == null) {
            handler = mCallbackHandler;
        }
        if (handler != null) {
            handler.post(new Runnable() {
                @Override
                public void run() {
                    callback.onRecvFailed(baseData, errMsg);
                }
            });
        } else {
            callback.onRecvFailed(baseData, errMsg);
        }
    }

    // 发送数据成功回调
    public void onSendSuccess(final BaseSocketCallback callback, final BaseData baseData) {
        Handler handler = callback.getCallbackHandler();
        if (handler == null) {
            handler = mCallbackHandler;
        }
        if (handler != null) {
            handler.post(new Runnable() {
                @Override
                public void run() {
                    callback.onSendSuccess(baseData);
                }
            });
        } else {
            callback.onSendSuccess(baseData);
        }
    }

    // 发送数据失败回调
    public void onSendFailed(final BaseSocketCallback callback, final BaseData baseData, final String errMsg) {
        Handler handler = callback.getCallbackHandler();
        if (handler == null) {
            handler = mCallbackHandler;
        }
        if (handler != null) {
            handler.post(new Runnable() {
                @Override
                public void run() {
                    callback.onSendFailed(baseData, errMsg);
                }
            });
        } else {
            callback.onSendFailed(baseData, errMsg);
        }
    }

    public BaseSocket() {
        init();
    }

    public BaseSocket(String ip, int port) {
        mIP = ip;
        mPort = port;
        init();
    }

    private void init() {
        if (mExit == true) {
            try {
                mExit = false;
                if (mTimeoutThread == null) {
                    mTimeoutThread = new HandlerThread("BaseSocketTimeoutThread");
                    mTimeoutThread.start();
                    mTimeoutHandler = new Handler(mTimeoutThread.getLooper());
                    // 网络异常重连
                    mTimeoutHandler.postDelayed(new Runnable() {
                        @Override
                        public void run() {
                            if (mExit == false && isConnected() == false) {
                                connect();
                            }
                            mTimeoutHandler.postDelayed(this, 3000);
                        }
                    },3000);

                    // 心跳
                    mTimeoutHandler.postDelayed(new Runnable() {
                        @Override
                        public void run() {
                            if (mExit == false && isConnected() == false) {
                                send(DataPackage.DATA_PROTO_CMD_BEATHEART, new BeatHeart(), new SocketCallback());
                            }
                            mTimeoutHandler.postDelayed(this, 120 * mSecond);
                        }
                    }, 120*mSecond);
                }
                sendThread.start();
                recvThread.start();
            } catch (Exception e) {
                Logger.d(TAG, e.toString());
                mExit = true;
            }
        }

        if (mCallbackMap == null) {
            mCallbackMap = new ConcurrentHashMap<Integer, BaseSocketCallback>();
        }
        if (mPushCallbackMap == null) {
            mPushCallbackMap = new ConcurrentHashMap<Integer, BaseSocketCallback>();
        }
    }

    public void exit() {
        if (mTimeoutThread != null) {
            mTimeoutThread.quitSafely();
        }
        mExit = true;
    }
    public void setSendTimeout(int timeout) {
        mSendTimeout = timeout;
    }

    public void setRecvTimeout(int timeout) { mRecvTimeout = timeout; }

    public boolean connect() {
        return connect(mIP, mPort);
    }

    public boolean connect(String ip, int port) {
        if (socket != null && socket.isConnected()) {
            this.close();
        }
        mIP = ip;
        mPort = port;

        try {
            socket = new Socket(ip, port);
            // 延迟40ms 影响不是很大
            socket.setTcpNoDelay(true);
            socket.setKeepAlive(true);
            in = new DataInputStream(socket.getInputStream());
            out = new DataOutputStream(socket.getOutputStream());
        } catch (IOException e) {
            Logger.i(TAG, "connect server failed");
            return false;
        }
        setConnected(true);
        Logger.d(TAG, "connect server success!!!!");
        return true;
    }
    public void setConnected(boolean connected) {
        synchronized (mConnectedLock) {
            mConnected = connected;
        }
    }
    public boolean isConnected() {
        /*
        try {
            if (socket == null)
                return false;

            socket.sendUrgentData(0);
            socket.sendUrgentData(0);
        } catch (Exception e) {
            return false;
        }
        */
        synchronized (mConnectedLock) {
            return mConnected;
        }
    }
    public void close() {
        try {
            if (socket != null) {
                socket.close();
                socket = null;
            }
            if (in != null) {
                in.close();
                in = null;
            }
            if (out != null) {
                out.close();
                out = null;
            }
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            setConnected(false);
        }
    }

    public static void sleep(int time) {
        try {
            Thread.sleep(time);
        } catch (Exception e) {
            Logger.d(TAG, e.toString());
        }
    }

    public Thread recvThread = new Thread() {
        @Override
        public void run() {
            short pkglen = 0;
            byte[] byteBuffer = new byte[DataPackage.MAX_DATA_PACKAGE_SIZE];
            while (mExit == false) {
                if (isConnected() == false) {
                    BaseSocket.sleep(1000);
                    continue;
                }
                try {
                    // 默认大端方式读取
                    pkglen = in.readShort();
                    if (pkglen > DataPackage.MAX_DATA_PACKAGE_SIZE || pkglen < BaseData.DATA_PACKAGE_HEAD_SIZE) {
                        Logger.i(TAG, "data package size error: " + pkglen);
                        continue;
                    }
                    in.readFully(byteBuffer, 2, pkglen - 2);
                } catch (IOException e) {
                    Logger.d(TAG, e.toString());
                    setConnected(false);
                    continue;
                }
                // 大端高位在低地址，低位在高地址
                byteBuffer[0] = (byte) ((pkglen >> 8) & 0xff);
                byteBuffer[1] = (byte) (pkglen & 0xff);

                BaseData baseData = DataPackage.convertBytesToObject(ByteBuffer.wrap(byteBuffer, 0, pkglen));
                if (baseData == null) {
                    Logger.d(TAG, "recv null data package");
                    continue;
                }
                BaseSocketCallback callback = mPushCallbackMap.get(baseData.getPushID());
                if (callback != null) {
                    onRecvSuccess(callback, baseData);
                    continue;
                }

                callback = mCallbackMap.remove(baseData.getID());
                if (callback == null) {
                    Logger.d(TAG, "recv callback is null");
                    continue;
                }
                onRecvSuccess(callback, baseData);
            }
        }
    };

    public void send(int cmd, final BaseData baseData, BaseSocketCallback callback) {
        baseData.command = cmd;
        putSendTask(baseData);
        mCallbackMap.put(baseData.getID(), callback);
        /*
        mTimeoutHandler.postDelayed(new Runnable() {
            @Override
            public void run() {
                BaseSocketCallback cb = mCallbackMap.get(baseData.getID());
                if (cb != null) {
                    onSendFailed(cb, baseData, "send timeout");
                }
            }
        }, mSendTimeout);
        */
        mTimeoutHandler.postDelayed(new Runnable() {
            @Override
            public void run() {
                BaseSocketCallback cb = mCallbackMap.remove(baseData.getID());
                if (cb != null) {
                    onRecvFailed(cb, baseData, "recv timeout");
                }
            }
        }, mRecvTimeout);
    }

    public Thread sendThread = new Thread() {
        @Override
        public void run() {
            boolean sleep = false;
            while (mExit == false) {
                if (isConnected() == false) {
                    BaseSocket.sleep(1000);
                    continue;
                }
                BaseData baseData = getSendTask();
                if (baseData == null) {
                    continue;
                }
                BaseSocketCallback callback = mCallbackMap.get(baseData.getID());
                if (callback == null) {
                    Logger.d(TAG, "callback is null, can not send basedata");
                    continue;
                }
                try {
                    out.write(DataPackage.convertObjectToBytes(baseData));
                } catch (Exception e) {
                    onSendFailed(callback,baseData, e.toString());
                    setConnected(false);
                    continue;
                }
                onSendSuccess(callback, baseData);
            }
        }
    };

    public BaseData getSendTask() {
        BaseData baseData = null;
        sendTaskLock.lock();
        if (sendTaskList.isEmpty()) {
            sendTaskLock.unlock();
            try {
                synchronized (sendTaskLockObject) {
                    sendTaskLockObject.wait();
                }
            }catch (Exception e) {
                Logger.d(TAG, e.toString());
            }
            return null;
        }
        try {
            baseData = sendTaskList.poll();
        } catch (Exception e) {
            Logger.d(TAG, e.toString());
        }
        sendTaskLock.unlock();
        return baseData;
    }

    public void putSendTask(BaseData baseData) {
        sendTaskLock.lock();
        if (sendTaskList.isEmpty()) {
            synchronized (sendTaskLockObject) {
                sendTaskLockObject.notify();
            }
        }
        sendTaskList.offer(baseData);
        sendTaskLock.unlock();
    }
}
