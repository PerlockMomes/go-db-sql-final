package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	parcelAdd, err := s.db.Exec(
		"insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, fmt.Errorf("произошла ошибка добавления новой записи в базу данных, %w\n", err)
	}
	number, err := parcelAdd.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("произошла ошибка возврата последнего добаленного номера, %w\n", err)
	}
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	parcel := Parcel{}
	getRow := s.db.QueryRow("select * from parcel where number = :number",
		sql.Named("number", number))
	err := getRow.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		return parcel, fmt.Errorf("посылка № %d не найдена, %w", number, err)
	}
	return parcel, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var sliceParcel []Parcel
	getRows, err := s.db.Query("select * from parcel where client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("клиент %d не найден, %w\n", client, err)
	}
	defer getRows.Close()
	for getRows.Next() {
		parcel := Parcel{}
		err := getRows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("клиент № %d, произошла ошибка копирования, %w\n", client, err)
		}
		sliceParcel = append(sliceParcel, parcel)
	}
	return sliceParcel, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("update parcel set status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("не получилось обновить статус посылки № %d, %w\n", number, err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("update parcel set address = :address where number = :number and status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return fmt.Errorf("не удалось обновить статус посылки № %d, %w\n", number, err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("delete from parcel where number = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return fmt.Errorf("не удалось удалить посылку № %d, %w\n", number, err)
	}
	return nil
}
