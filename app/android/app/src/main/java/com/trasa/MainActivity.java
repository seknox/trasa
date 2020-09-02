package com.trasa;

import com.facebook.react.ReactActivity;
import android.os.Bundle; // here
import com.zoontek.rnbootsplash.RNBootSplash; // <- add this necessary import

public class MainActivity extends ReactActivity {
  @Override
  protected void onCreate(Bundle savedInstanceState) {
    RNBootSplash.init(R.drawable.screen, MainActivity.this); // <- display the generated bootsplash.xml drawable over our MainActivity
    super.onCreate(savedInstanceState);
  }

  /**
   * Returns the name of the main component registered from JavaScript. This is used to schedule
   * rendering of the component.
   */
  @Override
  protected String getMainComponentName() {
    return "trasa";
  }
}
