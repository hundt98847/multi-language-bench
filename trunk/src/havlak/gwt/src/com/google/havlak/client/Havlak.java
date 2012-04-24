package com.google.havlak.client;

import com.google.havlak.shared.FieldVerifier;
import com.google.havlak.shared.cfg.BasicBlock;
import com.google.havlak.shared.cfg.BasicBlockEdge;
import com.google.havlak.shared.cfg.CFG;
import com.google.havlak.shared.lsg.LSG;
import com.google.havlak.shared.lsg.SimpleLoop;
import com.google.havlak.shared.havlakloopfinder.HavlakLoopFinder;
import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;

class LoopTesterApp {
  public LoopTesterApp() {
    cfg = new CFG();
    lsg = new LSG();
    root = cfg.createNode(0);
  }

  // Create 4 basic blocks, corresponding to and if/then/else clause
  // with a CFG that looks like a diamond
  public int buildDiamond(int start) {
    int bb0 = start;
    new BasicBlockEdge(cfg, bb0, bb0 + 1);
    new BasicBlockEdge(cfg, bb0, bb0 + 2);
    new BasicBlockEdge(cfg, bb0 + 1, bb0 + 3);
    new BasicBlockEdge(cfg, bb0 + 2, bb0 + 3);

    return bb0 + 3;
  }

  // Connect two existing nodes
  public void buildConnect(int start, int end) {
    new BasicBlockEdge(cfg, start, end);
  }

  // Form a straight connected sequence of n basic blocks
  public int buildStraight(int start, int n) {
    for (int i = 0; i < n; i++) {
      buildConnect(start + i, start + i + 1);
    }
    return start + n;
  }

  // Construct a simple loop with two diamonds in it
  public int buildBaseLoop(int from) {
    int header = buildStraight(from, 1);
    int diamond1 = buildDiamond(header);
    int d11 = buildStraight(diamond1, 1);
    int diamond2 = buildDiamond(d11);
    int footer = buildStraight(diamond2, 1);
    buildConnect(diamond2, d11);
    buildConnect(diamond1, header);

    buildConnect(footer, from);
    footer = buildStraight(footer, 1);
    return footer;
  }

  public  CFG        cfg;
  public  LSG        lsg;
  private BasicBlock root;
}

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class Havlak implements EntryPoint {
  /**
   * The message displayed to the user when the server cannot be reached or
   * returns an error.
   */
  private static final String SERVER_ERROR = "An error occurred while "
      + "attempting to contact the server. Please check your network "
      + "connection and try again.";

  /**
   * Create a remote service proxy to talk to the server-side Greeting service.
   */
  private final GreetingServiceAsync greetingService = GWT.create(GreetingService.class);

  /**
   * This is the entry point method.
   */
  public void onModuleLoad() {
    final Button sendButton = new Button("Send to Havlak");
    final TextBox nameField = new TextBox();
    nameField.setText("GWT User");
    final Label errorLabel = new Label();

    // We can add style names to widgets
    sendButton.addStyleName("sendButton");

    // Add the nameField and sendButton to the RootPanel
    // Use RootPanel.get() to get the entire body element
    RootPanel.get("nameFieldContainer").add(nameField);
    RootPanel.get("sendButtonContainer").add(sendButton);
    RootPanel.get("errorLabelContainer").add(errorLabel);

    // Focus the cursor on the name field when the app loads
    nameField.setFocus(true);
    nameField.selectAll();

    // Create the popup dialog box
    final DialogBox dialogBox = new DialogBox();
    dialogBox.setText("Loop Recognition");
    dialogBox.setAnimationEnabled(true);
    final Button closeButton = new Button("Close");
    // We can set the id of a widget by accessing its Element
    closeButton.getElement().setId("closeButton");
    final Label textToServerLabel = new Label();
    final HTML serverResponseLabel = new HTML();
    VerticalPanel dialogVPanel = new VerticalPanel();
    dialogVPanel.addStyleName("dialogVPanel");
    dialogVPanel.add(textToServerLabel);
    dialogVPanel.add(new HTML("<br><b>Algorithm replies:</b>"));
    dialogVPanel.add(serverResponseLabel);
    dialogVPanel.setHorizontalAlignment(VerticalPanel.ALIGN_RIGHT);
    dialogVPanel.add(closeButton);
    dialogBox.setWidget(dialogVPanel);

    // Add a handler to close the DialogBox
    closeButton.addClickHandler(new ClickHandler() {
      public void onClick(ClickEvent event) {
        dialogBox.hide();
        sendButton.setEnabled(true);
        sendButton.setFocus(true);
      }
    });

    // Create a handler for the sendButton and nameField
    class MyHandler implements ClickHandler, KeyUpHandler {
      /**
       * Fired when the user clicks on the sendButton.
       */
      public void onClick(ClickEvent event) {
        sendNameToServer();
      }

      /**
       * Fired when the user types in the nameField.
       */
      public void onKeyUp(KeyUpEvent event) {
        if (event.getNativeKeyCode() == KeyCodes.KEY_ENTER) {
          sendNameToServer();
        }
      }

      /**
       * Send the name from the nameField to the server and wait for a response.
       */
      private void sendNameToServer() {
        // First, we validate the input.
        errorLabel.setText("");
        String textToServer = nameField.getText();

        //=======================================================
        // HAVLAK
        //=======================================================
        long start = System.currentTimeMillis();
        String result = "Welcome to LoopTesterApp, GWT edition<br>";
        LoopTesterApp app = new LoopTesterApp();
        app.cfg.createNode(0);
        app.lsg.dump();
        app.buildBaseLoop(0);
        app.cfg.createNode(1);
        new BasicBlockEdge(app.cfg, 0, 2);
        int found = 0;

        result += "15000 dummy loops<br>";
        serverResponseLabel.setHTML(result);
        for (int dummyloop = 0; dummyloop < 1; dummyloop++) {
          HavlakLoopFinder finder = new HavlakLoopFinder(app.cfg, app.lsg);
          finder.findLoops();
        }

        result += "Constructing CFG...<br>";
        serverResponseLabel.setHTML(result);
        int n = 2;

        for (int parlooptrees = 0; parlooptrees < 10; parlooptrees++) {
          app.cfg.createNode(n + 1);
          app.buildConnect(2, n + 1);
          n = n + 1;

          for (int i = 0; i < 2; i++) {
            int top = n;
            n = app.buildStraight(n, 1);
            for (int j = 0; j < 25; j++) {
              n = app.buildBaseLoop(n);
            }
            int bottom = app.buildStraight(n, 1);
            app.buildConnect(n, top);
            n = bottom;
          }
          app.buildConnect(n, 1);
        }

        result += "Performing Loop Recognition\n1 Iteration<br>";
        HavlakLoopFinder finder = new HavlakLoopFinder(app.cfg, app.lsg);
        finder.findLoops();
        long t = System.currentTimeMillis() - start;
        result += "Found: " + app.lsg.getNumLoops() + " in " +
            t + " [ms]";

        result += "Another 100 iterations...<br>";
        for (int i = 0; i < 100; i++) {
          HavlakLoopFinder finder2 = new HavlakLoopFinder(app.cfg, new LSG());
          finder2.findLoops();
        }
        t = System.currentTimeMillis() - start;
        result += "<br>Found: " + app.lsg.getNumLoops() + " in " +
            t + " [ms]";

        //=======================================================

        sendButton.setEnabled(false);
        dialogBox.setText("Find Loops");
        serverResponseLabel.setHTML(result);
        dialogBox.center();
        closeButton.setFocus(true);

        /*
        //        if (!FieldVerifier.isValidName(textToServer)) {
        //          errorLabel.setText("Please enter at least four characters");
        //          return;
        //        }
        
        // Then, we send the input to the server.
        textToServerLabel.setText(textToServer);
        serverResponseLabel.setText("");
        greetingService.greetServer(textToServer, new AsyncCallback<String>() {
          public void onFailure(Throwable caught) {
            // Show the RPC error message to the user
            dialogBox.setText("Remote Procedure Call - Failure");
            serverResponseLabel.addStyleName("serverResponseLabelError");
            serverResponseLabel.setHTML(SERVER_ERROR);
            dialogBox.center();
            closeButton.setFocus(true);
          }

          public void onSuccess(String result) {
            dialogBox.setText("Finding Loops...");
            serverResponseLabel.removeStyleName("serverResponseLabelError");

            serverResponseLabel.setHTML(result);
            dialogBox.center();
            closeButton.setFocus(true);
          }
        });
        */
      }
    }

    // Add a handler to send the name to the server
    MyHandler handler = new MyHandler();
    sendButton.addClickHandler(handler);
    nameField.addKeyUpHandler(handler);
  }
}
