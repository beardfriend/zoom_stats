import React, { useState, useEffect } from 'react';
import DateDrawer from './components/date';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js';
import { subDays, format } from 'date-fns';
import { Bar } from 'react-chartjs-2';
import './App.css';
import axios from 'axios';
ChartJS.register(
  CategoryScale,
  LinearScale,
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
  const [data2, setData2] = useState({ labels: [], datasets: [] });
  // 화제의 채팅
  const [popularChat, setPopularChat] = useState([]);
  const [option, setOption] = useState({});
  const [option2, setOption2] = useState({});

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
          console.log(jsonData);
          setPopularChat(jsonData);
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };
    PopuplarChatFetch();
  }, [date]);

  useEffect(() => {
    const fetchData2 = async () => {
      try {
        const response = await axios.get(
          'http://localhost:8080/api/chats/most-reactors',
          {
            params: {
              start_date_time: '2023-06-06T00:00:00',
              end_date_time: '2023-06-06T23:59:59'
            }
          }
        );
        const jsonData = response.data;
        let labels = [];
        let dda = [];
        jsonData.forEach((d) => {
          labels.push(d.name);
          dda.push(d.count);
        });

        const dd = {
          labels: labels,
          datasets: [
            {
              label: '좋아요 개수',
              data: dda,
              backgroundColor: 'red',
              borderWidth: 1
            }
          ]
        };
        setData2(dd);
        const options = {
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
        };
        setOption2(options);
      } catch (error) {
        console.error('Error:', error);
      }
    };

    fetchData2();
  }, []);

  return (
    <div className='App'>
      <>
        <DateDrawer date={date} setDate={setDate} />
        <h1>이어드림 줌 통계</h1>
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
            <Bar data={data2} options={option2} />
          </div>
        </div>
        <div>
          <h3>화제의 채팅 TOP10</h3>
          <div
            style={{ textAlign: 'left', margin: `0 auto`, maxWidth: '1200px' }}
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
      </>
    </div>
  );
}

export default App;
