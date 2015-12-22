package com.example;

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;

public class Socket extends Thread{
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
            socket.setSoTimeout(10000);
            in = new DataInputStream(socket.getInputStream());
            out = new DataOutputStream(socket.getOutputStream());
        } catch (Exception e) {
            e.printStackTrace();
            return false;
        }
        return true;
    }

    public boolean send(DataPackage dataPackage) {
        byte[] data = dataPackage.pack();
        try {
            out.write(data);
        } catch (IOException e) {
            e.printStackTrace();
            return false;
        }
        return true;
    }

    public DataPackage recv() {
        short pkglen = 0;
        pkglen = 0;
        try {
            in.readFully(byteBuffer, 0, 2);
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }

        pkglen += (int) byteBuffer[1];
        pkglen += (int) byteBuffer[0] << 8;

        if (pkglen > MAX_DATA_PACKAGE_SIZE || pkglen < DataPackage.DATA_PACKAGE_HEAD_SIZE) {
            System.out.println("data package size error");
            System.out.printf("%d %d %d%n", pkglen, byteBuffer[0], byteBuffer[1]);
            return null;
        }

        try {
            in.readFully(byteBuffer, 2, pkglen - 2);
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
        public short version;
        public int command;
        public int errorCode;
        public final static int DATA_PACKAGE_HEAD_SIZE = 12;
        public byte[] packageBody = null;

        public void unpack(byte[] b, int offset, int len) {
            ByteBuffer byteBuffer = ByteBuffer.wrap(b, offset, len);
            byteBuffer.order(ByteOrder.BIG_ENDIAN);
            pacakgeLen = byteBuffer.getShort();
            version = byteBuffer.getShort();
            command = byteBuffer.getInt();
            errorCode = byteBuffer.getInt();
            if (pacakgeLen > DATA_PACKAGE_HEAD_SIZE) {
                packageBody = new byte[ pacakgeLen - DATA_PACKAGE_HEAD_SIZE];
                byteBuffer.get(packageBody);
            } else {
                packageBody = null;
            }
        }

        public byte[] pack() {
            ByteBuffer byteBuffer = ByteBuffer.allocate(pacakgeLen);
            byteBuffer.order(ByteOrder.BIG_ENDIAN);
            byteBuffer.putShort(pacakgeLen);
            byteBuffer.putShort(version);
            byteBuffer.putInt(command);
            byteBuffer.putInt(errorCode);
            if (packageBody != null) {
                byteBuffer.put(packageBody);
            }
            return byteBuffer.array();
        }

        public String toString() {
            String body = "null";
            if (packageBody != null) {
                body = new String(packageBody);
            }
            return String.format("package size: %d%ncommand: %d%nversion: %d%nnerrorCode: %d%n" +
                            "DataBody: %s%n", pacakgeLen,command, version, errorCode, body);
        }
    }
}
