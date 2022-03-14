import React from "react";
import { Link } from "react-router-dom";
import SidebarOption from "./SidebarOption";
import HomeIcon from "@material-ui/icons/Home";
import PermIdentityIcon from "@material-ui/icons/PermIdentity";
import SettingsIcon from "@material-ui/icons/Settings";
import ExitToAppIcon from "@material-ui/icons/ExitToApp";

import "./Sidebar.css";

function Sidebar() {
  return (
    <div className="mt-5">
      <Link to="/home">
        <SidebarOption Icon={HomeIcon} text="Home" />
      </Link>
      <Link to="/mypage">
        <SidebarOption Icon={PermIdentityIcon} text="My Page" />
      </Link>
      <Link to="/settings">
        <SidebarOption Icon={SettingsIcon} text="Settings" />
      </Link>
      <Link to="/logout">
        <SidebarOption Icon={ExitToAppIcon} text="Logout" />
      </Link>
    </div>
  );
}

export default Sidebar;
