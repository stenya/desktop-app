import { Notification } from "electron";
import { PauseStateEnum } from "@/store/types";
import store from "@/store";

var connectionNotification = null;
var handlerOnNotificationClick = null;

export function InitNotifications(onNotificationClickHandler) {
  handlerOnNotificationClick = onNotificationClickHandler;
}

store.subscribe((mutation) => {
  try {
    switch (mutation.type) {
      case "vpnState/connectionInfo":
      case "vpnState/disconnected":
      case "vpnState/pauseState":
        {
          const isConnected = store.getters["vpnState/isConnected"];
          let messageText = "";

          // first notification should be only 'connected'
          if (isConnected !== true && !connectionNotification) break;

          // initialize notification text (body)
          if (isConnected === true)
            messageText = `Connected\n${getConnectedInfoText()}`;
          else messageText = "Disconnected";
          if (store.state.vpnState.pauseState === PauseStateEnum.Paused)
            messageText = "Paused";

          // if notification already shown and have the same body - do nothing (needed notification already shown)
          if (
            connectionNotification &&
            messageText == connectionNotification.body
          )
            break;

          // if notification object not created yet - create it (it is a first time notification shown)
          if (!connectionNotification) {
            connectionNotification = new Notification({
              title: "IVPN",
              body: messageText,
            });
          }
          connectionNotification.body = messageText;

          // set on-click handler
          connectionNotification.on("click", () => {
            if (handlerOnNotificationClick) {
              try {
                handlerOnNotificationClick();
              } catch (e) {
                console.error("Notification click handler error: ", e);
              }
            }
          });

          connectionNotification.show();
        }
        break;
    }
  } catch (e) {
    console.error("Error in store subscriber (notifications):", e);
  }
});

function getConnectedInfoText() {
  function text(svr) {
    if (svr == null) return "";
    return `${svr.city}, ${svr.country_code}`;
  }

  try {
    if (!store.getters["vpnState/isConnected"]) return "";
    let ret = text(store.state.settings.serverEntry);
    if (store.state.settings.isMultiHop)
      ret = `${ret} -> ${text(store.state.settings.serverExit)}`;

    return ret;
  } catch (e) {
    console.error(e);
    return "";
  }
}
