package com.example;


import net.sf.json.JSONObject;

/**
 * Created by wuxiangan on 2015/12/22.
 */
public class Client {
    public String username;
    public String password;
    private  Socket socket;

    public Client(String username, String password) {
        this.username = username;
        this.password = password;
        this.socket = new Socket();
    }

    public Client(String username, String password, String svrIP, int svrPort) {
        this.username = username;
        this.password = password;
        this.socket = new Socket(svrIP, svrPort);
    }

    public void close() {
        this.socket.close();
    }

    public boolean isConnect() {
        return this.socket.isConnected();
    }

    public boolean connect(String svrIP, int svrPort) {
        return this.socket.connect(svrIP, svrPort);
    }

    public boolean connect() {
        return this.socket.connect();
    }

    public boolean heartbeat() {
        Socket.DataPackage dataPackage = new Socket().new DataPackage();
        dataPackage.pacakgeLen = (short)dataPackage.DATA_PACKAGE_HEAD_SIZE;
        dataPackage.command = 10000;

        if (false == this.socket.send(dataPackage)) {
            return false;
        }

        dataPackage = this.socket.recv();
        if (dataPackage == null || dataPackage.errorCode != 0) {
            return false;
        }

        return true;
    }
    public boolean login() {
        Socket.DataPackage dataPackage = new Socket().new DataPackage();
        JSONObject jsonObject = new JSONObject();
        jsonObject.put("uid", username);
        jsonObject.put("password", password);

        dataPackage.packageBody = jsonObject.toString().getBytes();
        dataPackage.command = 10002;
        dataPackage.pacakgeLen = (short)(dataPackage.DATA_PACKAGE_HEAD_SIZE + dataPackage.packageBody.length);

        if (false == this.socket.send(dataPackage)) {
            return false;
        }

        dataPackage = this.socket.recv();
        if (dataPackage == null || dataPackage.errorCode != 0) {
            return false;
        }

        return true;
    }
}
