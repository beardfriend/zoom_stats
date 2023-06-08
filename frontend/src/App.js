import axios from 'axios';
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  Title,
  Tooltip
} from 'chart.js';
import { format, subDays } from 'date-fns';
import React, { useEffect, useState } from 'react';
import { Bar, Line } from 'react-chartjs-2';
import './App.css';
import DateDrawer from './components/date';
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend
);
function App() {
  const [date, setDate] = useState(subDays(new Date(), 1));
  // 가장 따봉을 많이 받은 사람
  const [DDabong, setDDabong] = useState({ labels: [], datasets: [] });
  const [DDabongOption, _] = useState({
    responsive: true,
    plugins: {
      legend: {
        position: 'top'
      },
      title: {
        display: true,
        text: '공감 많이 받은 TOP 10'
      }
    }
  });

  // 공감많이 누른 사람
  const [reactor, setReactor] = useState({ labels: [], datasets: [] });
  const [reactorOption, setReactorOption] = useState({
    responsive: true,
    plugins: {
      legend: {
        position: 'top'
      },
      title: {
        display: true,
        text: '공감 많이 한 TOP 10'
      }
    }
  });
  // 월별 일간 참여도
  const [participation, setParticipation] = useState({
    labels: [],
    datasets: []
  });
  const [participationOption, setParticipationOption] = useState({
    responsive: true,
    plugins: {
      legend: {
        position: 'top'
      },
      title: {
        display: true,
        text: '월별 일간 참여도'
      }
    }
  });
  // 화제의 채팅
  const [popularChat, setPopularChat] = useState([]);

  useEffect(() => {
    const DDabongFetch = async () => {
      try {
        const response = await axios.get(
          'http://localhost:8080/api/chats/most-reacted-people',
          {
            params: {
              start_date_time: `${format(date, 'yyyy-MM-dd')}T00:00:00`,
              end_date_time: `${format(date, 'yyyy-MM-dd')}T23:59:59`
            }
          }
        );
        const jsonData = response.data;

        if (jsonData === null) {
          setDDabong({ labels: [], datasets: [] });
          return;
        }
        let labels = [];
        let data = [];
        jsonData.forEach((d) => {
          labels.push(d.name);
          data.push(d.count);
        });

        const dd = {
          labels: labels,
          datasets: [
            {
              label: '좋아요 개수',
              data: data,
              backgroundColor: 'rgba(0, 123, 255, 0.5)',
              borderColor: 'rgba(0, 123, 255, 1)',
              borderWidth: 1
            }
          ]
        };
        setDDabong(dd);
      } catch (error) {
        console.error('Error:', error);
      }
    };
    DDabongFetch();
  }, [date]);

  useEffect(() => {
    const PopuplarChatFetch = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/chats/', {
          params: {
            start_date_time: `${format(date, 'yyyy-MM-dd')}T00:00:00`,
            end_date_time: `${format(date, 'yyyy-MM-dd')}T23:59:59`
          }
        });

        const jsonData = response.data;

        if (jsonData === null) {
          setPopularChat([]);
          return;
        } else {
          setPopularChat(jsonData);
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };
    PopuplarChatFetch();
  }, [date]);

  useEffect(() => {
    const ParticipationFetch = async () => {
      try {
        const beforeDate = subDays(date, 30);

        const response = await axios.get(
          'http://localhost:8080/api/chats/participation',
          {
            params: {
              start_date_time: `${format(beforeDate, 'yyyy-MM-dd')}T00:00:00`,
              end_date_time: `${format(date, 'yyyy-MM-dd')}T23:59:59`
            }
          }
        );

        const jsonData = response.data;
        if (jsonData === null) {
          setParticipation({ labels: [], datasets: [] });
          return;
        }

        let labels = [];
        let data = [];
        jsonData.forEach((d) => {
          labels.push(d.date);
          data.push(d.count);
        });

        const dd = {
          labels: labels,
          datasets: [
            {
              label: '참여 카운트',
              data: data,
              borderColor: 'rgb(53, 162, 235)',
              backgroundColor: 'rgba(53, 162, 235, 0.5)'
            }
          ]
        };
        setParticipation(dd);
      } catch (error) {
        console.error('Error:', error);
      }
    };
    ParticipationFetch();
  }, [date]);

  useEffect(() => {
    const reactorFetch = async () => {
      try {
        const response = await axios.get(
          'http://localhost:8080/api/chats/most-reactors',
          {
            params: {
              start_date_time: `${format(date, 'yyyy-MM-dd')}T00:00:00`,
              end_date_time: `${format(date, 'yyyy-MM-dd')}T23:59:59`
            }
          }
        );
        const jsonData = response.data;
        if (jsonData === null) {
          setReactor({ labels: [], datasets: [] });
          return;
        }
        const labels = [];
        const data = [];
        jsonData.forEach((d) => {
          labels.push(d.name);
          data.push(d.count);
        });

        const datas = {
          labels: labels,
          datasets: [
            {
              label: '좋아요 개수',
              data: data,
              backgroundColor: 'red',
              borderWidth: 1
            }
          ]
        };

        setReactor(datas);
      } catch (error) {
        console.error('Error:', error);
      }
    };

    reactorFetch();
  }, [date]);

  return (
    <div
      className='App'
      style={{ padding: '0 10rem', boxSizing: 'content-box' }}
    >
      <>
        <h1>이어드림 줌 통계</h1>
        <DateDrawer date={date} setDate={setDate} />

        <div style={{ display: 'grid', gridTemplateColumns: `1fr 1fr` }}>
          <div>
            <h3>공감 많이 받은</h3>
            {DDabong.datasets.legnth === 0 ? (
              <h1>empty</h1>
            ) : (
              <Bar data={DDabong} options={DDabongOption} />
            )}
          </div>
          <div>
            <h3>공감 많이 한</h3>
            <Bar data={reactor} options={reactorOption} />
          </div>

          <div>
            <h3>월별 일간 참여도</h3>
            <Line data={participation} options={participationOption} />
          </div>
          <div>
            <h3>화제의 채팅 TOP10</h3>
            <div
              style={{
                textAlign: 'left',
                margin: `0 auto`,
                maxWidth: '1200px'
              }}
            >
              {popularChat.map((d, index) => (
                <div key={index} className='chatItem'>
                  <p>
                    {d.chatted_at} {d.sender} : {d.text}{' '}
                    <span style={{ fontWeight: 'bold' }}>
                      / {d.react_ids.length}❤️
                    </span>
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </>
    </div>
  );
}

export default App;
