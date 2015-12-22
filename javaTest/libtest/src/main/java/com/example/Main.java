package com.example;

/**
 * Created by wuxiangan on 2015/12/18.
 */
public class Main {
    public static void main(String []args) {
        //Client client = new Client("100000", "test", "192.168.20.51", 9101);
        Client client = new Client("1000000", "test", "127.0.0.1", 9100);

        if (false == client.connect()) {
            System.out.println("connect server failed");
            return;
        }

        if (false == client.heartbeat()) {
            System.out.println("client send heartbeat failed");
            return;
        }
        if (false == client.login()) {
            System.out.println("login failed");
            return ;
        }

        System.out.println("login success");
        return;
    }
}
