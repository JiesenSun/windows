package ui.base;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.os.Bundle;
import android.util.AttributeSet;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import com.example.wuxiangan.bangbang.R;

import java.util.concurrent.ConcurrentHashMap;

import util.ActivityCollector;

/**
 * Created by wuxiangan on 2016/1/15.
 */
public class TitleLayout extends LinearLayout{
    public Button titleLeft;
    public Button titleRight;
    public TextView middleText;
    private static Activity activity=null;
    public TitleLayout(Context context, AttributeSet attributeSet) {
        super(context, attributeSet);
        LayoutInflater.from(context).inflate(R.layout.title_layout, this);

        titleLeft = (Button)findViewById(R.id.title_left);
        titleRight = (Button)findViewById(R.id.title_right);
        middleText = (TextView)findViewById(R.id.title_text);
        titleLeft.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                TitleLayout.activity = (Activity)getContext();
                if (ActivityCollector.activityNum() == 2) {
                    AlertDialog.Builder dialog = new AlertDialog.Builder(TitleLayout.activity);
                    dialog.setTitle(R.string.exit_tip);
                    dialog.setCancelable(false);
                    dialog.setPositiveButton(R.string.OK, new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            ActivityCollector.finishAll();
                        }
                    });
                    dialog.setNegativeButton(R.string.CANCEL,null);
                    dialog.show();
                }
            }
        });

        titleRight.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                Toast.makeText(getContext(), "You clicked edit button", Toast.LENGTH_SHORT).show();
            }
        });
    }

    public void setTitleLeft(CharSequence text) {
        titleLeft.setText(text);
    }

    public void setTitleRight(CharSequence text) {
        titleRight.setText(text);
    }

    public void setTitle(CharSequence text){
        middleText.setText(text);
    }
}
