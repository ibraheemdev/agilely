import React from "react"
import { machineSelectors } from "@redux/machine"
import { useSelector } from "react-redux";

const Status = () => {
  const status = useSelector(state => machineSelectors.getStatus(state))
  return (
    <div>
      {status}
    </div>
  )
}
export default Status