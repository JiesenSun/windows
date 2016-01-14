package client;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.annotations.Expose;

import org.json.JSONObject;

import java.nio.ByteBuffer;

import protocol.BaseData;
import socket.Socket;
import socket.SocketCallback;
import util.Logger;

/**
 * Created by wuxiangan on 2016/1/5.
 */
public class Login extends BaseData{
    private static final String TAG = "Login";
    //public String mUserName;
    @Expose
    public long mUserName;
    @Expose
    public String mPassword;

    public  long mSessionID;

    public void convertObjectToBytes(ByteBuffer byteBuffer) {
        JSONObject jsonObject = new JSONObject();
        try {
            jsonObject.put("uid", mUserName);
            jsonObject.put("password", mPassword);
        } catch (Exception e) {
            Logger.d(TAG, e.toString());
        }
        byteBuffer.put(jsonObject.toString().getBytes());
    }

    public BaseData convertBytesToObject(ByteBuffer byteBuffer) {
        int dataSize = byteBuffer.limit() - byteBuffer.position();
        byte[] dataBuf = new byte[dataSize];
        byteBuffer.get(dataBuf);
        try {
            JSONObject jsonObject = new JSONObject(new String(dataBuf));
            mUserName = jsonObject.getLong("uid");
            mSessionID = jsonObject.getLong("sid");
        } catch (Exception e) {
            Logger.d(TAG, e.toString());
        }
        return this;
    }

    public void login() {
        Gson gson = new GsonBuilder().excludeFieldsWithoutExposeAnnotation().create();
        this.mUserName = 1000000;
        this.mPassword = "test";
        String josnString = gson.toJson(this);
        Logger.d(TAG, gson.toJson(this));
        Socket.send(10002, this, new SocketCallback());
    }
}
