package com.example;

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;

public class Socket extends Thread{
    private static int MIN_DATA_PACKAGE_SIZE = 24;
    private static int MAX_DATA_PACKAGE_SIZE = 2048;
    private String serverIP;
    private int serverPort;
    private java.net.Socket socket;
    private DataInputStream in;
    private DataOutputStream out;
    private byte[] byteBuffer = new byte[MAX_DATA_PACKAGE_SIZE];
    private DataPackagHandle dataHandle = new DataPackagHandle() {
        @Override
        public boolean handle(DataPackage dataPackage) {
            System.out.println(dataPackage);
            return true;
        }
    };

    public Socket() {}
    public Socket(String serverIP, int serverPort) {
        this.serverIP = serverIP;
        this.serverPort = serverPort;
    }

    public void setDataHandle(DataPackagHandle handle) {
        dataHandle = handle;
    }

    public void close() {
        try {
            if (socket != null) {
                socket.close();
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
        socket = null;
    }

    public boolean connect() {
        return connect(serverIP, serverPort);
    }

    public boolean connect(String serverIP, int serverPort) {
        this.serverIP = serverIP;
        this.serverPort = serverPort;
        try {
            socket = new java.net.Socket(serverIP, serverPort);
            socket.setTcpNoDelay(true);
            in = new DataInputStream(socket.getInputStream());
            out = new DataOutputStream(socket.getOutputStream());
        } catch (Exception e) {
            e.printStackTrace();
            return false;
        }
        return true;
    }

    public void send(DataPackage dataPackage) {
        byte[] data = dataPackage.pack();
        try {
            out.write(data);
        } catch (IOException e) {
            e.printStackTrace();
        }
        System.out.println("send data package len:");
        System.out.println(data.length);
    }

    public DataPackage recv() {
        short pkglen = 0;
        pkglen = 0;
        try {
            in.readFully(byteBuffer, 0, 4);
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }

        pkglen += (int) byteBuffer[1];
        pkglen += (int) byteBuffer[0] << 8;
        //pkglen += (int)byteBuffer[2] << 16;
        //pkglen += (int)byteBuffer[3] << 24;

        if (pkglen > MAX_DATA_PACKAGE_SIZE || pkglen < MIN_DATA_PACKAGE_SIZE) {
            System.out.println("data package size error");
            return null;
        }

        try {
            in.readFully(byteBuffer, 4, pkglen - 4);
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }

        DataPackage dataPackage = new DataPackage();
        dataPackage.unpack(byteBuffer, 0, pkglen);
        return dataPackage;
    }

    public boolean isConnected() {
        return socket != null;
    }


    public void run() {
        while (socket != null) {
            DataPackage dataPackage = recv();
            if (dataPackage == null) {
                close();
                return;
            }
            if (dataHandle != null && false == dataHandle.handle(dataPackage)) {
                //xxx
            }
        }
    }


    public interface DataPackagHandle {
        public boolean handle(DataPackage dataPackage);
    }

    public class DataPackage {
        public short pacakgeLen;
        public short command;
        public short version;
        public short sequence;
        public int sessionID;
        public long userID;
        public int errorCode;
        public int bodyLen;
        public final static int DATA_PACKAGE_HEAD_SIZE = 24;
        public byte[] packageBody = new byte[MAX_DATA_PACKAGE_SIZE];

        public void unpack(byte[] b, int offset, int len) {
            ByteBuffer byteBuffer = ByteBuffer.wrap(b, offset, len);
            byteBuffer.order(ByteOrder.BIG_ENDIAN);
            pacakgeLen = byteBuffer.getShort();
            command = byteBuffer.getShort();
            version = byteBuffer.getShort();
            sequence = byteBuffer.getShort();
            sessionID = byteBuffer.getInt();
            userID = byteBuffer.getLong();
            errorCode = byteBuffer.getInt();
            bodyLen = pacakgeLen - DATA_PACKAGE_HEAD_SIZE;
            byteBuffer.get(packageBody, 0, (int) bodyLen);
        }

        public byte[] pack() {
            ByteBuffer byteBuffer = ByteBuffer.allocate(MAX_DATA_PACKAGE_SIZE);
            byteBuffer.order(ByteOrder.BIG_ENDIAN);
            byteBuffer.putShort(pacakgeLen);
            byteBuffer.putShort(command);
            byteBuffer.putShort(version);
            byteBuffer.putShort(sequence);
            byteBuffer.putInt(sessionID);
            byteBuffer.putLong(userID);
            byteBuffer.putInt(errorCode);
            byteBuffer.put(packageBody, 0, bodyLen);
            return byteBuffer.array();
        }

        public String toString() {
            String body = new String(packageBody, 0, bodyLen);
            return String.format("package size: %d%ncommand: %d%nversion: %d%nsequence: %d%n" +
                            "sessionID: %d%nuserID: %d%nerrorCode: %d%nDataBody: %s%n", pacakgeLen,
                    command, version, sequence, sessionID, userID, errorCode, body);
        }
    }
}
