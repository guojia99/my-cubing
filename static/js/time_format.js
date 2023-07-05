/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/5 下午9:43.
 *  * Author: guojia(https://github.com/guojia99)
 */


function formatTimeByProject(number, project) {
    if (typeof number !== "number" || isNaN(number) || number < 0) {
        return "Invalid input";
    }

    // Check if the project is '最少步'
    const isMinSteps = project === "最少步";

    // Convert seconds to minutes, seconds, and milliseconds
    const minutes = Math.floor(number / 60);
    const seconds = Math.floor(number % 60);
    const milliseconds = Math.floor((number % 1) * 1000);

    // Format minutes and seconds with leading zeros if necessary
    const formattedMinutes = isMinSteps ? minutes.toString() : String(minutes).padStart(2, "0");
    const formattedSeconds = minutes > 0 ? String(seconds).padStart(2, "0") : seconds.toString();

    // Check if there are non-zero minutes or seconds
    const hasMinutes = minutes > 0;
    const hasSeconds = seconds > 0 || (isMinSteps && number > 0);

    // Output the result based on the conditions
    if (isMinSteps && Number.isInteger(number)) {
        if (number === 0){
            return "DNF"
        }
        return String(number);
    } else if (hasMinutes && hasSeconds) {
        return `${formattedMinutes}:${formattedSeconds}.${String(milliseconds).padStart(3, "0")}`;
    } else if (hasMinutes && !hasSeconds) {
        return `${formattedMinutes}.${String(milliseconds).padStart(3, "0")}`;
    } else if (!hasMinutes && hasSeconds) {
        return `${formattedSeconds}.${String(milliseconds).padStart(3, "0")}`;
    } else {
        return "DNF";
    }
}