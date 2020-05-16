import {Component, OnInit} from '@angular/core';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {formatDate} from '@angular/common';

interface Info {
  id: number;
  date: string;
  startTime: number;
  stopTime: number;
  breakStartTime: number;
  breakStopTime: number;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})

export class AppComponent implements OnInit {
  title = 'app';
  infos: Info[];
  todayInfo: Info;

  editDate: string;
  editStartTime: string;
  editStopTime: string;
  editBreakStartTime: string;
  editBreakStopTime: string;

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    this.updateState();
  }

  updateState() {
    this.today();
    this.getInfo();
  }

  getInfo() {
    {
      this.http.get<Info[]>('http://localhost:5300/api/info').subscribe(data => {
          for (const item of data) {
            console.log('id', +item.id);
            console.log('date', +item.date);
            console.log('startTime', +item.startTime);
            console.log('stopTime', +item.stopTime);
            console.log('breakStartTime', +item.breakStartTime);
            console.log('breakStopTime', +item.breakStopTime);
          }
          this.infos = data;
        },
        (err: HttpErrorResponse) => {
          if (err.error instanceof Error) {
            console.log('Client-side error occured.');
          } else {
            console.log('Server-side error occured.');
          }
        }
      );
    }
  }

  today() {
    return this.http.get<Info>('http://localhost:5300/api/today').subscribe(
      data => {
        console.log('data', data);
        this.todayInfo = data;
      });
  }

  startButtonIsDisplayed() {
    return this.todayInfo == null;
  }

  start() {
    return this.http.post('http://localhost:5300/api/start', null).subscribe(
      () => this.updateState(),
      data => {
        console.log('data', data);
      });
  }

  startBreakButtonIsDisplayed() {
    return this.todayInfo != null && this.todayInfo.breakStartTime === 0;
  }

  startBreak() {
    return this.http.post('http://localhost:5300/api/start-break', null).subscribe(
      () => this.updateState(),
      data => {
        console.log('data', data);
      });
  }

  stopBreakButtonIsDisplayed() {
    return this.todayInfo != null && this.todayInfo.breakStartTime !== 0 && this.todayInfo.breakStopTime === 0;
  }

  stopBreak() {
    return this.http.post('http://localhost:5300/api/stop-break', null).subscribe(
      () => this.updateState(),
      data => {
        console.log('data', data);
      });
  }

  stopButtonIsDisplayed() {
    return this.todayInfo != null &&
      this.todayInfo.startTime !== 0 &&
      this.todayInfo.breakStartTime !== 0 &&
      this.todayInfo.breakStopTime !== 0;
  }

  stop() {
    return this.http.post('http://localhost:5300/api/stop', null).subscribe(
      () => this.updateState(),
      data => {
        console.log('data', data);
      });
  }

  edit(infoEdit: Info) {
    this.editDate = infoEdit.date;
    this.editStartTime = formatDate(
      infoEdit.startTime * 1000, 'HH:mm:ss', 'en-US', ''
    ).toString();
    this.editStopTime = formatDate(
      infoEdit.stopTime * 1000, 'HH:mm:ss', 'en-US', ''
    ).toString();
  }

  editBreak(infoEdit: Info) {
    this.editDate = infoEdit.date;
    this.editBreakStartTime = formatDate(
      infoEdit.breakStartTime * 1000, 'HH:mm:ss', 'en-US', ''
    ).toString();
    this.editBreakStopTime = formatDate(
      infoEdit.breakStopTime * 1000, 'HH:mm:ss', 'en-US', ''
    ).toString();
  }

  editSave() {
    this.http.post('http://localhost:5300/api/edit', null, {
      params: {
        date: this.editDate,
        startTime: (new Date(this.editDate + ' ' + this.editStartTime).getTime() / 1000).toString(),
        stopTime: (new Date(this.editDate + ' ' + this.editStopTime).getTime() / 1000).toString(),
      }
    }).subscribe(
      () => this.getInfo(),
      data => {
        console.log('data', data);
      });
  }

  editBreakSave() {
    this.http.post('http://localhost:5300/api/edit-break', null, {
      params: {
        date: this.editDate,
        breakStartTime: (new Date(this.editDate + ' ' + this.editBreakStartTime).getTime() / 1000).toString(),
        breakStopTime: (new Date(this.editDate + ' ' + this.editBreakStopTime).getTime() / 1000).toString(),
      }
    }).subscribe(
      () => this.getInfo(),
      data => {
        console.log('data', data);
      });
  }
}
