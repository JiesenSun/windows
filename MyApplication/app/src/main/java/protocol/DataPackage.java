package protocol;

import android.util.Log;

import java.nio.ByteBuffer;
import java.util.InputMismatchException;
import java.util.concurrent.ConcurrentHashMap;

import client.BeatHeart;
import client.LoginT;
import util.Logger;

/**
 * Created by wuxiangan on 2015/12/31.
 */
public class DataPackage {
    public static final int DATA_PROTO_CMD_BEATHEART = 10000;
    public static final String TAG = "DataPackage";
    public static final int MAX_DATA_PACKAGE_SIZE = 2048;
    private static ConcurrentHashMap<Integer, Class<?>> mDataProtoMap = new ConcurrentHashMap<Integer, Class<?>>();

    public static void register() {
        registerDataProto(10000, BeatHeart.class);
        registerDataProto(10002, LoginT.class);

    }
    public static void registerDataProto(int cmd, Class<?> dataProto) {
        if (null == mDataProtoMap.get(cmd)) {
            mDataProtoMap.put(cmd, dataProto);
        }
    }

    public static BaseData get(int cmd) {
        BaseData baseData = null;
        try {
            Class<?> classLoader = mDataProtoMap.get(cmd);
            baseData = (BaseData)classLoader.newInstance();
        } catch (Exception e) {
            Logger.d(TAG, e.toString());
            return null;
        }
        return baseData;
    }

    public static BaseData convertBytesToObject(ByteBuffer byteBuffer) {
        BaseData dataPackage = null;
        try {
            BaseData baseData = new BaseData().convertBytesToObject(byteBuffer);
            dataPackage = get(baseData.command);
            dataPackage = dataPackage.convertBytesToObject(byteBuffer);
            dataPackage.initBaseData(baseData);
        } catch (Exception e) {
            Log.i(TAG, e.toString());
            return null;
        }

        return dataPackage;
    }

    public static byte[] convertObjectToBytes(BaseData baseData) {
        byte[] byteBuff = null;
        try {
            ByteBuffer byteBuffer = ByteBuffer.allocate(MAX_DATA_PACKAGE_SIZE);
            baseData.cloneBaseData().convertObjectToBytes(byteBuffer);
            baseData.convertObjectToBytes(byteBuffer);
            int position = byteBuffer.position();
            byteBuffer.position(0);
            byteBuffer.putShort((short)position);
            byteBuffer.position(position);
            byteBuffer.flip();
            byteBuff = new byte[position];
            byteBuffer.get(byteBuff, 0, position);
        } catch (Exception e) {
            Log.i(TAG, e.toString());
            return null;
        }
        return byteBuff;
    }

    public static byte[] convertObjectToBytes(BaseData baseData, int cmd) {
        baseData.command = cmd;
        return DataPackage.convertObjectToBytes(baseData);
    }
}
