import React, { useState } from 'react';
import { format, addDays, subDays } from 'date-fns';

const DateDrawer = ({ date, setDate }) => {
  const handleNextDay = () => {
    const nextDay = addDays(date, 1);
    setDate(nextDay);
  };

  const handlePreviousDay = () => {
    const previousDay = subDays(date, 1);
    setDate(previousDay);
  };

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center'
      }}
    >
      <button style={previousButtonStyle} onClick={handlePreviousDay}>
        &lt; Previous
      </button>
      <h2 style={{ margin: '0 20px' }}>{format(date, 'yyyy-MM-dd')}</h2>
      <button style={nextButtonStyle} onClick={handleNextDay}>
        Next &gt;
      </button>
    </div>
  );
};
const previousButtonStyle = {
  backgroundColor: '#4CAF50',
  border: 'none',
  color: 'white',
  padding: '10px 20px',
  textAlign: 'center',
  textDecoration: 'none',
  display: 'inline-block',
  fontSize: '16px',
  margin: '4px 2px',
  cursor: 'pointer',
  borderRadius: '4px'
};

const nextButtonStyle = {
  backgroundColor: 'white',
  border: '1px solid #4CAF50',
  color: '#4CAF50',
  padding: '10px 20px',
  textAlign: 'center',
  textDecoration: 'none',
  display: 'inline-block',
  fontSize: '16px',
  margin: '4px 2px',
  cursor: 'pointer',
  borderRadius: '4px'
};
export default DateDrawer;
