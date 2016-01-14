package client;

import java.nio.ByteBuffer;

import protocol.BaseData;

/**
 * Created by wuxiangan on 2016/1/5.
 */
public class BeatHeart extends BaseData {
    @Override
    public BaseData convertBytesToObject(ByteBuffer byteBuffer) {
        return this;
    }

    @Override
    public void convertObjectToBytes(ByteBuffer byteBuffer) {
    }
}
