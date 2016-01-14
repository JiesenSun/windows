package protocol;


import java.nio.ByteBuffer;

/**
 * Created by wuxiangan on 2015/12/31.
 */
public class BaseData {
    private static int mBaseDataID = 0;
    private static final Object mBaseDataIDLock = new Object();
    public final static int DATA_PACKAGE_HEAD_SIZE = 16;  // 包头大小

    private short packageLen;                             // 包长度
    private short version;                                // 版本
    public int command;                                   // 命令
    private int sequenceID;                               // 包ID
    private int errorCode;                                // 错误码

    public BaseData() {
        sequenceID = BaseData.BaseDataID();
    }
    public static int BaseDataID() {
        synchronized (mBaseDataIDLock) {
            if (mBaseDataID == Integer.MAX_VALUE) {
                mBaseDataID = 0;
            }
            return mBaseDataID++;
        }
    }

    public BaseData initBaseData(BaseData baseData) {
        this.packageLen = baseData.packageLen;
        this.version = baseData.version;
        this.command = baseData.command;
        this.sequenceID = baseData.sequenceID;
        this.errorCode = baseData.errorCode;
        return this;
    }

    public int getID() {
        return sequenceID;
    }
    public int getPushID() {
        return command;
    }
    public BaseData cloneBaseData() {
         return new BaseData().initBaseData(this);
    }
    // 子类重写此方法后无须再调用父类此方法，否则出错重复写
    public BaseData convertBytesToObject(ByteBuffer byteBuffer) {
        this.packageLen = byteBuffer.getShort();
        this.version = byteBuffer.getShort();
        this.command = byteBuffer.getInt();
        this.sequenceID = byteBuffer.getInt();
        this.errorCode = byteBuffer.getInt();
        return this;
    }

    public void convertObjectToBytes(ByteBuffer byteBuffer) {
        byteBuffer.putShort(this.packageLen);
        byteBuffer.putShort(this.version);
        byteBuffer.putInt(this.command);
        byteBuffer.putInt(this.sequenceID);
        byteBuffer.putInt(this.errorCode);
    }
}
