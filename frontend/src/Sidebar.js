import React from "react";
import { Link } from "react-router-dom";
import { Icon } from "@iconify/react";
import catIcon from "@iconify-icons/mdi/cat";
import SidebarOption from "./SidebarOption";
import HomeIcon from "@material-ui/icons/Home";
import NotificationsNoneIcon from "@material-ui/icons/NotificationsNone";
import PermIdentityIcon from "@material-ui/icons/PermIdentity";
import SettingsIcon from "@material-ui/icons/Settings";

import "./Sidebar.css";

function Sidebar() {
  return (
    <div>
      <Icon icon={catIcon} width="3rem" height="3rem" />
      <Link to="/home">
        <SidebarOption Icon={HomeIcon} text="Home" />
      </Link>
      {/* <SidebarOption Icon={NotificationsNoneIcon} text="Notifications" /> */}
      <Link to="/mypage">
        <SidebarOption Icon={PermIdentityIcon} text="My Page" />
      </Link>
      <Link to="/settings">
        <SidebarOption Icon={SettingsIcon} text="Settings" />
      </Link>
    </div>
  );
}

export default Sidebar;
