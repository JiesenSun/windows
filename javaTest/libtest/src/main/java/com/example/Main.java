package com.example;

/**
 * Created by wuxiangan on 2015/12/18.
 */
public class Main {
    public static void main(String []args) {
        Socket client = new Socket();

        if (false == client.connect("192.168.20.51", 9100)) {
            System.out.println("connect server failed!!!");
            return;
        }
        System.out.println("connect server success...");
        Socket.DataPackage dataPackage = client.recv();
        System.out.println(dataPackage);
        client.start();
        client.send(dataPackage);
        try {
            Thread.sleep(1000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        client.close();
    }
}
